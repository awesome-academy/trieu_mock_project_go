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

func (r *SkillRepository) FindByID(db *gorm.DB, id uint) (*models.Skill, error) {
	var skill models.Skill
	result := db.First(&skill, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &skill, nil
}

func (r *SkillRepository) SearchSkills(db *gorm.DB, limit, offset int) ([]models.Skill, int64, error) {
	var skills []models.Skill
	query := db.Model(&models.Skill{})

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&skills).Error; err != nil {
		return nil, 0, err
	}

	return skills, count, nil
}

func (r *SkillRepository) Create(db *gorm.DB, skill *models.Skill) error {
	return db.Create(skill).Error
}

func (r *SkillRepository) Update(db *gorm.DB, skill *models.Skill) error {
	return db.Save(skill).Error
}

func (r *SkillRepository) Delete(db *gorm.DB, id uint) error {
	return db.Delete(&models.Skill{}, id).Error
}

func (r *SkillRepository) ExistsUsersWithSkillID(db *gorm.DB, skillID uint) (bool, error) {
	var userSkill models.UserSkill
	err := db.
		Select("user_id").
		Where("skill_id = ?", skillID).
		First(&userSkill).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
