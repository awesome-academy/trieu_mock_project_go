package services

import (
	"context"
	"strings"
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
		Name:         strings.TrimSpace(req.Name),
		Abbreviation: strings.TrimSpace(req.Abbreviation),
	}

	err := s.positionRepository.Create(s.db.WithContext(c), position)
	if err != nil && err == gorm.ErrDuplicatedKey {
		return appErrors.ErrPositionAlreadyExists
	}
	return err
}

func (s *PositionService) UpdatePosition(c context.Context, id uint, req dtos.CreateOrUpdatePositionRequest) error {
	currentPosition, err := s.positionRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrNotFound
		}

		return err
	}

	currentPosition.Name = strings.TrimSpace(req.Name)
	currentPosition.Abbreviation = strings.TrimSpace(req.Abbreviation)

	err = s.positionRepository.Update(s.db.WithContext(c), currentPosition)
	if err != nil && err == gorm.ErrDuplicatedKey {
		return appErrors.ErrPositionAlreadyExists
	}
	return err
}

func (s *PositionService) DeletePosition(c context.Context, id uint) error {
	_, err := s.positionRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrNotFound
		}
		return err
	}

	existedUserUsePosition, err := s.positionRepository.ExistsUsersWithPositionID(s.db.WithContext(c), id)
	if err != nil {
		return err
	}
	if existedUserUsePosition {
		return appErrors.ErrPositionInUse
	}

	return s.positionRepository.Delete(s.db.WithContext(c), id)
}
