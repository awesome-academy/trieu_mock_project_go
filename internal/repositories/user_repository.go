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

func (r *UserRepository) SearchUsers(db *gorm.DB, name *string, teamId *uint, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	query := db.Model(&models.User{})

	result := query

	if name != nil {
		result = result.Where("name LIKE ?", "%"+*name+"%")
	}

	if teamId != nil {
		result = result.Where("current_team_id = ?", *teamId)
	}

	var count int64
	result = result.Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = result.
		Preload("CurrentTeam").
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

func (r *UserRepository) UpdateUser(db *gorm.DB, user *models.User) error {
	return db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"name":            user.Name,
			"email":           user.Email,
			"birthday":        user.Birthday,
			"current_team_id": user.CurrentTeamID,
			"position_id":     user.PositionID,
			"role":            user.Role,
		}).Error
}

func (r *UserRepository) UpdateUserSkills(db *gorm.DB, userID uint, skills []models.UserSkill) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Delete existing skills
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserSkill{}).Error; err != nil {
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

func (r *UserRepository) UpdateUsersCurrentTeamToNullByTeamID(db *gorm.DB, teamID uint) error {
	return db.Model(&models.User{}).
		Where("current_team_id = ?", teamID).
		Update("current_team_id", nil).Error
}
