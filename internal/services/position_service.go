package services

import (
	"context"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/repositories"

	"gorm.io/gorm"
)

type PositionService struct {
	db                 *gorm.DB
	positionRepository *repositories.PositionRepository
}

func NewPositionService(db *gorm.DB, positionRepository *repositories.PositionRepository) *PositionService {
	return &PositionService{db: db, positionRepository: positionRepository}
}

func (s *PositionService) GetAllPositionsSummary(c context.Context) []dtos.PositionSummary {
	positions, err := s.positionRepository.FindAllPositionsSummary(s.db.WithContext(c))
	if err != nil {
		return []dtos.PositionSummary{}
	}

	positionDtos := make([]dtos.PositionSummary, 0, len(positions))
	for _, position := range positions {
		positionDtos = append(positionDtos, dtos.PositionSummary{
			ID:   position.ID,
			Name: position.Name,
		})
	}
	return positionDtos
}
