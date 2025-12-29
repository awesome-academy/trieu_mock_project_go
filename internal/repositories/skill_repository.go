package repositories

import (
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type SkillRepository struct {
}

func NewSkillRepository() *SkillRepository {
	return &SkillRepository{}
}

func (r *SkillRepository) FindAllSkillSummary(db *gorm.DB) ([]models.Skill, error) {
	var skills []models.Skill
	result := db.
		Select("id", "name").
		Find(&skills)
	if result.Error != nil {
		return nil, result.Error
	}
	return skills, nil
}
