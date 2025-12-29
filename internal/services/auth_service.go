package services

import (
	"context"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/models"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type AuthService struct {
	db   *gorm.DB
	repo *repositories.UserRepository
}

func NewAuthService(db *gorm.DB, repo *repositories.UserRepository) *AuthService {
	return &AuthService{db: db, repo: repo}
}

func (s *AuthService) Login(c context.Context, email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(s.db.WithContext(c), email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, appErrors.ErrNotFound
	}

	if !s.VerifyPassword(password, user.Password) {
		return nil, appErrors.ErrInvalidCredentials
	}

	return user, nil
}

func (s *AuthService) VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
