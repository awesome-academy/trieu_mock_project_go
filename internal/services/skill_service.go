package services

import (
	"context"
	"strings"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type SkillService struct {
	db              *gorm.DB
	skillRepository *repositories.SkillRepository
}

func NewSkillService(db *gorm.DB, skillRepository *repositories.SkillRepository) *SkillService {
	return &SkillService{db: db, skillRepository: skillRepository}
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
		return nil, err
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
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}

	return helpers.MapSkillToSkillSummary(skill), nil
}

func (s *SkillService) CreateSkill(c context.Context, req dtos.CreateOrUpdateSkillRequest) error {
	skill := &models.Skill{
		Name: strings.TrimSpace(req.Name),
	}

	err := s.skillRepository.Create(s.db.WithContext(c), skill)

	if err != nil && appErrors.IsDuplicatedEntryError(err) {
		return appErrors.ErrSkillAlreadyExists
	}
	return err
}

func (s *SkillService) UpdateSkill(c context.Context, id uint, req dtos.CreateOrUpdateSkillRequest) error {
	currentSkill, err := s.skillRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrNotFound
		}

		return err
	}

	currentSkill.Name = strings.TrimSpace(req.Name)

	err = s.skillRepository.Update(s.db.WithContext(c), currentSkill)
	if err != nil && appErrors.IsDuplicatedEntryError(err) {
		return appErrors.ErrSkillAlreadyExists
	}
	return err
}

func (s *SkillService) DeleteSkill(c context.Context, id uint) error {
	_, err := s.skillRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrNotFound
		}
		return err
	}

	existedUserUseSkill, err := s.skillRepository.ExistsUsersWithSkillID(s.db.WithContext(c), id)
	if err != nil {
		return err
	}
	if existedUserUseSkill {
		return appErrors.ErrSkillInUse
	}

	return s.skillRepository.Delete(s.db.WithContext(c), id)
}
