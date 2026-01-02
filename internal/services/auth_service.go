package services

import (
	"context"
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
}

func NewAuthService(db *gorm.DB, repo *repositories.UserRepository, activityLogService *ActivityLogService) *AuthService {
	return &AuthService{db: db, repo: repo, activityLogService: activityLogService}
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

func (s *AuthService) VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
