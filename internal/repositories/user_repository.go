package repositories

import (
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) FindByID(db *gorm.DB, id uint) (*models.User, error) {
	var user models.User
	result := db.
		Preload("CurrentTeam").
		Preload("Position").
		Preload("Projects").
		Preload("Skills").
		First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
