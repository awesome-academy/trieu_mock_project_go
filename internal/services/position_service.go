package services

import (
	"context"
	"fmt"
	"strings"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/internal/types"
	"trieu_mock_project_go/internal/utils"
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type PositionService struct {
	db                 *gorm.DB
	positionRepository *repositories.PositionRepository
	activityLogService *ActivityLogService
}

func NewPositionService(db *gorm.DB, positionRepository *repositories.PositionRepository, activityLogService *ActivityLogService) *PositionService {
	return &PositionService{db: db, positionRepository: positionRepository, activityLogService: activityLogService}
}

func (s *PositionService) GetAllPositionsSummary(c context.Context) []dtos.PositionSummary {
	positions, err := s.positionRepository.FindAllPositionsSummary(s.db.WithContext(c))
	if err != nil {
		return []dtos.PositionSummary{}
	}

	return helpers.MapPositionsToPositionSummaries(positions)
}

func (s *PositionService) SearchPositions(c context.Context, limit, offset int) (*dtos.PositionSearchResponse, error) {
	positions, totalCount, err := s.positionRepository.SearchPositions(s.db.WithContext(c), limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	return &dtos.PositionSearchResponse{
		Positions: helpers.MapPositionsToPositionDtos(positions),
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
			return nil, appErrors.ErrPositionNotFound
		}
		return nil, appErrors.ErrInternalServerError
	}

	return helpers.MapPositionToPositionDto(position), nil
}

func (s *PositionService) CreatePosition(c context.Context, req dtos.CreateOrUpdatePositionRequest) error {
	position := &models.Position{
		Name:         strings.TrimSpace(req.Name),
		Abbreviation: strings.TrimSpace(req.Abbreviation),
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.positionRepository.Create(tx, position); err != nil {
			if appErrors.IsDuplicatedEntryError(err) {
				return appErrors.ErrPositionAlreadyExists
			}
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.CreatePosition, position.ID, position.Name); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *PositionService) UpdatePosition(c context.Context, id uint, req dtos.CreateOrUpdatePositionRequest) error {
	currentPosition, err := s.positionRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrPositionNotFound
		}

		return appErrors.ErrInternalServerError
	}

	currentPosition.Name = strings.TrimSpace(req.Name)
	currentPosition.Abbreviation = strings.TrimSpace(req.Abbreviation)

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.positionRepository.Update(tx, currentPosition); err != nil {
			if appErrors.IsDuplicatedEntryError(err) {
				return appErrors.ErrPositionAlreadyExists
			}
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.UpdatePosition, currentPosition.ID, currentPosition.Name); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *PositionService) DeletePosition(c context.Context, id uint) error {
	position, err := s.positionRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrPositionNotFound
		}
		return appErrors.ErrInternalServerError
	}

	existedUserUsePosition, err := s.positionRepository.ExistsUsersWithPositionID(s.db.WithContext(c), id)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	if existedUserUsePosition {
		return appErrors.ErrPositionInUse
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.positionRepository.Delete(tx, id); err != nil {
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.DeletePosition, position.ID, position.Name); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *PositionService) ExportPositionsToCSV(c context.Context) ([][]string, error) {
	positions, err := s.positionRepository.FindAllPositionsSummary(s.db.WithContext(c))
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	data := [][]string{{"ID", "Name", "Abbreviation"}}
	for _, p := range positions {
		data = append(data, []string{
			fmt.Sprintf("%d", p.ID),
			p.Name,
			p.Abbreviation,
		})
	}
	return data, nil
}

func (s *PositionService) ImportPositionsFromCSV(c context.Context, data [][]string) error {
	if len(data) <= 1 {
		return appErrors.ErrNoCSVDataToImport
	}

	positionsMap := make(map[string]*dtos.PositionImportData)
	uniqueNames := utils.NewSet[string]()

	for i, row := range data {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			return fmt.Errorf("row %d: invalid number of columns", i+1)
		}

		name := strings.TrimSpace(row[0])
		abbreviation := strings.TrimSpace(row[1])

		if name == "" || abbreviation == "" {
			return fmt.Errorf("row %d: name and abbreviation are required", i+1)
		}

		if !uniqueNames.Has(name) {
			uniqueNames.Add(name)
			positionsMap[name] = &dtos.PositionImportData{
				Name:         name,
				Abbreviation: abbreviation,
			}
		}
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		positions := make([]models.Position, 0, len(positionsMap))
		for _, name := range uniqueNames.ToSlice() {
			p := positionsMap[name]
			positions = append(positions, models.Position{
				Name:         p.Name,
				Abbreviation: p.Abbreviation,
			})
		}

		if err := s.positionRepository.CreateInBatches(tx, positions, 100); err != nil {
			if appErrors.IsDuplicatedEntryError(err) {
				return appErrors.ErrPositionAlreadyExists
			}
			return appErrors.ErrInternalServerError
		}

		activityLogsToInsert := make([]models.ActivityLog, 0, len(positions))
		for _, p := range positions {
			activityLog, err := s.activityLogService.createLogActivityModel(c, types.CreatePosition, p.ID, p.Name)
			if err != nil {
				return err
			}
			activityLogsToInsert = append(activityLogsToInsert, *activityLog)
		}

		if err := s.activityLogService.createInBatches(tx, activityLogsToInsert, 100); err != nil {
			return appErrors.ErrInternalServerError
		}

		return nil
	})
}
