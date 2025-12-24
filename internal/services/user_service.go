package services

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type UserService struct {
	db   *gorm.DB
	repo *repositories.UserRepository
}

func NewUserService(db *gorm.DB, repo *repositories.UserRepository) *UserService {
	return &UserService{db: db, repo: repo}
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(s.db, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if !s.VerifyPassword(password, user.Password) {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *UserService) VerifyPassword(plainPassword, hashedPassword string) bool {
	hash := sha256.Sum256([]byte(plainPassword))
	hashString := fmt.Sprintf("%x", hash)
	return hashString == hashedPassword
}
