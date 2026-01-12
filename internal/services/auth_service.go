package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/internal/types"
	"trieu_mock_project_go/models"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type AuthService struct {
	db                 *gorm.DB
	repo               *repositories.UserRepository
	activityLogService *ActivityLogService
	redisService       *RedisService
}

func NewAuthService(db *gorm.DB, repo *repositories.UserRepository, activityLogService *ActivityLogService, redisService *RedisService) *AuthService {
	return &AuthService{db: db, repo: repo, activityLogService: activityLogService, redisService: redisService}
}

func (s *AuthService) CreateWSTicket(ctx context.Context, userID uint, email string) (string, *appErrors.AppError) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", appErrors.ErrInternalServerError
	}
	ticket := hex.EncodeToString(b)

	key := fmt.Sprintf("ws_ticket:%s", ticket)
	value := fmt.Sprintf("%d:%s", userID, email)

	// Save to redis with 1 minute expiration
	err := s.redisService.Set(ctx, key, value, time.Minute)
	if err != nil {
		return "", appErrors.ErrInternalServerError
	}

	return ticket, nil
}

func (s *AuthService) ConsumeWSTicket(ctx context.Context, ticket string) (uint, string, *appErrors.AppError) {
	key := fmt.Sprintf("ws_ticket:%s", ticket)
	value, err := s.redisService.Get(ctx, key)
	if err != nil {
		return 0, "", appErrors.ErrInvalidToken
	}

	_ = s.redisService.Del(ctx, key)

	var userID uint
	var email string
	_, err = fmt.Sscanf(value, "%d:%s", &userID, &email)
	if err != nil {
		return 0, "", appErrors.ErrInternalServerError
	}

	return userID, email, nil
}

func (s *AuthService) Login(c context.Context, email, password string, isAdmin bool) (*models.User, error) {
	user, err := s.repo.FindByEmail(s.db.WithContext(c), email)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	if user == nil {
		return nil, appErrors.ErrUserNotFound
	}

	if !s.VerifyPassword(password, user.Password) {
		return nil, appErrors.ErrInvalidCredentials
	}

	var activityType types.ActivityLog
	if isAdmin {
		activityType = types.AdminSignIn
	} else {
		activityType = types.UserSignIn
	}
	s.activityLogService.LogActivity(c, activityType, user.ID, user.Email)

	return user, nil
}

func (s *AuthService) Logout(c context.Context, userID uint, token, email string) *appErrors.AppError {
	err := s.DeleteToken(c, userID, token)
	if err != nil {
		return appErrors.ErrInternalServerError
	}

	return s.activityLogService.LogActivity(c, types.UserSignOut, userID, email)
}

func (s *AuthService) StoreToken(ctx context.Context, userID uint, token string, expiration time.Duration) error {
	tokenHash := s.hashToken(token)
	key := fmt.Sprintf("token:%d:%s", userID, tokenHash)
	return s.redisService.Set(ctx, key, "true", expiration)
}

func (s *AuthService) IsTokenStoreValid(ctx context.Context, userID uint, token string) (bool, error) {
	tokenHash := s.hashToken(token)
	key := fmt.Sprintf("token:%d:%s", userID, tokenHash)
	return s.redisService.Exists(ctx, key)
}

func (s *AuthService) DeleteToken(ctx context.Context, userID uint, token string) error {
	tokenHash := s.hashToken(token)
	key := fmt.Sprintf("token:%d:%s", userID, tokenHash)
	return s.redisService.Del(ctx, key)
}

func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

func (s *AuthService) VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
