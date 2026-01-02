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
		Preload("UserSkill.Skill").
		First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) SearchUsers(db *gorm.DB, teamId *uint, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	query := db.Model(&models.User{})

	result := query

	if teamId != nil {
		result = result.Where("current_team_id = ?", *teamId)
	}

	var count int64
	result = result.Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = result.
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	return users, count, nil
}

func (r *UserRepository) CreateUser(db *gorm.DB, user *models.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreateUserSkills(db *gorm.DB, userSkills []models.UserSkill) error {
	if len(userSkills) > 0 {
		if err := db.Create(&userSkills).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) UpdateUser(db *gorm.DB, user *models.User, skills []models.UserSkill) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Update user basic info
		if err := tx.Model(user).Updates(map[string]interface{}{
			"name":            user.Name,
			"email":           user.Email,
			"birthday":        user.Birthday,
			"position_id":     user.PositionID,
			"current_team_id": user.CurrentTeamID,
		}).Error; err != nil {
			return err
		}

		// Delete existing skills
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.UserSkill{}).Error; err != nil {
			return err
		}

		// Insert new skills
		if len(skills) > 0 {
			if err := tx.Create(&skills).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
