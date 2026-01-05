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

type SkillService struct {
	db                 *gorm.DB
	skillRepository    *repositories.SkillRepository
	activityLogService *ActivityLogService
}

func NewSkillService(db *gorm.DB, skillRepository *repositories.SkillRepository, activityLogService *ActivityLogService) *SkillService {
	return &SkillService{db: db, skillRepository: skillRepository, activityLogService: activityLogService}
}

func (s *SkillService) GetAllSkillsSummary(c context.Context) []dtos.SkillSummary {
	skills, err := s.skillRepository.FindAllSkillSummary(s.db.WithContext(c))
	if err != nil {
		return []dtos.SkillSummary{}
	}

	return helpers.MapSkillsToSkillSummaries(skills)
}

func (s *SkillService) SearchSkills(c context.Context, limit, offset int) (*dtos.SkillSearchResponse, error) {
	skills, totalCount, err := s.skillRepository.SearchSkills(s.db.WithContext(c), limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	return &dtos.SkillSearchResponse{
		Skills: helpers.MapSkillsToSkillSummaries(skills),
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}, nil
}

func (s *SkillService) GetSkillByID(c context.Context, id uint) (*dtos.SkillSummary, error) {
	skill, err := s.skillRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.ErrSkillNotFound
		}
		return nil, appErrors.ErrInternalServerError
	}

	return helpers.MapSkillToSkillSummary(skill), nil
}

func (s *SkillService) CreateSkill(c context.Context, req dtos.CreateOrUpdateSkillRequest) error {
	skill := &models.Skill{
		Name: strings.TrimSpace(req.Name),
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.skillRepository.Create(tx, skill); err != nil {
			if appErrors.IsDuplicatedEntryError(err) {
				return appErrors.ErrSkillAlreadyExists
			}
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.CreateSkill, skill.ID, skill.Name); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *SkillService) UpdateSkill(c context.Context, id uint, req dtos.CreateOrUpdateSkillRequest) error {
	currentSkill, err := s.skillRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrSkillNotFound
		}
		return appErrors.ErrInternalServerError
	}

	currentSkill.Name = strings.TrimSpace(req.Name)

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.skillRepository.Update(tx, currentSkill); err != nil {
			if appErrors.IsDuplicatedEntryError(err) {
				return appErrors.ErrSkillAlreadyExists
			}
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.UpdateSkill, currentSkill.ID, currentSkill.Name); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *SkillService) DeleteSkill(c context.Context, id uint) error {
	skill, err := s.skillRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrSkillNotFound
		}
		return appErrors.ErrInternalServerError
	}

	existedUserUseSkill, err := s.skillRepository.ExistsUsersWithSkillID(s.db.WithContext(c), id)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	if existedUserUseSkill {
		return appErrors.ErrSkillInUse
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.skillRepository.Delete(tx, id); err != nil {
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.DeleteSkill, skill.ID, skill.Name); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *SkillService) ExportSkillsToCSV(c context.Context) ([][]string, error) {
	skills, err := s.skillRepository.FindAllSkillSummary(s.db.WithContext(c))
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	data := [][]string{{"ID", "Name"}}
	for _, s := range skills {
		data = append(data, []string{
			fmt.Sprintf("%d", s.ID),
			s.Name,
		})
	}
	return data, nil
}

func (s *SkillService) ImportSkillsFromCSV(c context.Context, data [][]string) error {
	if len(data) <= 1 {
		return appErrors.ErrNoCSVDataToImport
	}

	uniqueNames := utils.NewSet[string]()

	for i, row := range data {
		if i == 0 {
			continue
		}
		if len(row) < 1 {
			return fmt.Errorf("row %d: invalid number of columns", i+1)
		}

		name := strings.TrimSpace(row[0])
		if name == "" {
			return fmt.Errorf("row %d: name is required", i+1)
		}

		if !uniqueNames.Has(name) {
			uniqueNames.Add(name)
		}
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		skills := make([]models.Skill, 0, uniqueNames.Size())
		for _, name := range uniqueNames.ToSlice() {
			skills = append(skills, models.Skill{
				Name: name,
			})
		}

		if err := s.skillRepository.CreateInBatches(tx, skills, 100); err != nil {
			if appErrors.IsDuplicatedEntryError(err) {
				return appErrors.ErrSkillAlreadyExists
			}
			return appErrors.ErrInternalServerError
		}

		activityLogsToInsert := make([]models.ActivityLog, 0, len(skills))
		for _, skill := range skills {
			activityLog, err := s.activityLogService.createLogActivityModel(c, types.CreateSkill, skill.ID, skill.Name)
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
