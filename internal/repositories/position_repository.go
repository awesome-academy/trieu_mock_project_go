package repositories

import (
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type PositionRepository struct {
}

func NewPositionRepository() *PositionRepository {
	return &PositionRepository{}
}

func (r *PositionRepository) FindByName(db *gorm.DB, name string) ([]models.Position, error) {
	var positions []models.Position
	result := db.
		Where("name = ?", name).
		Find(&positions)
	if result.Error != nil {
		return nil, result.Error
	}
	return positions, nil
}

func (r *PositionRepository) FindAllPositionsSummary(db *gorm.DB) ([]models.Position, error) {
	var positions []models.Position
	result := db.Find(&positions)
	if result.Error != nil {
		return nil, result.Error
	}
	return positions, nil
}

func (r *PositionRepository) FindByID(db *gorm.DB, id uint) (*models.Position, error) {
	var position models.Position
	result := db.First(&position, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &position, nil
}

func (r *PositionRepository) SearchPositions(db *gorm.DB, limit, offset int) ([]models.Position, int64, error) {
	var positions []models.Position
	query := db.Model(&models.Position{})

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&positions).Error; err != nil {
		return nil, 0, err
	}

	return positions, count, nil
}

func (r *PositionRepository) Create(db *gorm.DB, position *models.Position) error {
	return db.Create(position).Error
}

func (r *PositionRepository) Update(db *gorm.DB, position *models.Position) error {
	return db.Save(position).Error
}

func (r *PositionRepository) Delete(db *gorm.DB, id uint) error {
	return db.Delete(&models.Position{}, id).Error
}
