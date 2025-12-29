package services

import (
	"context"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/models"

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

func (s *PositionService) SearchPositions(c context.Context, limit, offset int) (*dtos.PositionSearchResponse, error) {
	positions, totalCount, err := s.positionRepository.SearchPositions(s.db.WithContext(c), limit, offset)
	if err != nil {
		return nil, err
	}

	positionDtos := make([]dtos.Position, 0, len(positions))
	for _, p := range positions {
		positionDtos = append(positionDtos, dtos.Position{
			ID:           p.ID,
			Name:         p.Name,
			Abbreviation: p.Abbreviation,
		})
	}

	return &dtos.PositionSearchResponse{
		Positions: positionDtos,
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}, nil
}

func (s *PositionService) GetPositionByID(c context.Context, id uint) (*dtos.Position, error) {
	position, err := s.positionRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}

	return &dtos.Position{
		ID:           position.ID,
		Name:         position.Name,
		Abbreviation: position.Abbreviation,
	}, nil
}

func (s *PositionService) CreatePosition(c context.Context, req dtos.CreateOrUpdatePositionRequest) error {
	position := &models.Position{
		Name:         req.Name,
		Abbreviation: req.Abbreviation,
	}

	existeds, err := s.positionRepository.FindByName(s.db.WithContext(c), req.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if len(existeds) > 0 {
		return appErrors.ErrPositionAlreadyExists
	}

	return s.positionRepository.Create(s.db.WithContext(c), position)
}

func (s *PositionService) UpdatePosition(c context.Context, id uint, req dtos.CreateOrUpdatePositionRequest) error {
	currentPosition, err := s.positionRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrNotFound
		}

		return err
	}

	if currentPosition.Name != req.Name {
		existeds, err := s.positionRepository.FindByName(s.db.WithContext(c), req.Name)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if len(existeds) > 0 {
			return appErrors.ErrPositionAlreadyExists
		}

		currentPosition.Name = req.Name
		currentPosition.Abbreviation = req.Abbreviation

		return s.positionRepository.Update(s.db.WithContext(c), currentPosition)
	}
	return nil
}

func (s *PositionService) DeletePosition(c context.Context, id uint) error {
	_, err := s.positionRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrNotFound
		}
		return err
	}

	return s.positionRepository.Delete(s.db.WithContext(c), id)
}
