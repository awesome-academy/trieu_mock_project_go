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

func (r *PositionRepository) FindAllPositionsSummary(db *gorm.DB) ([]models.Position, error) {
	var positions []models.Position
	result := db.Find(&positions)
	if result.Error != nil {
		return nil, result.Error
	}
	return positions, nil
}
