package services

import (
	"context"
	"strings"
	"time"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type ProjectService struct {
	db                *gorm.DB
	projectRepository *repositories.ProjectRepository
}

func NewProjectService(db *gorm.DB, projectRepository *repositories.ProjectRepository) *ProjectService {
	return &ProjectService{db: db, projectRepository: projectRepository}
}

func (s *ProjectService) GetAllProjectSummary(c context.Context) []dtos.ProjectSummary {
	projects, err := s.projectRepository.FindAllProjectSummary(s.db.WithContext(c))
	if err != nil {
		return []dtos.ProjectSummary{}
	}

	return helpers.MapProjectsToProjectSummaries(projects)
}

func (s *ProjectService) SearchProjects(c context.Context, teamID *uint, limit, offset int) (*dtos.ProjectSearchResponse, error) {
	projects, totalCount, err := s.projectRepository.SearchProjects(s.db.WithContext(c), teamID, limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	return &dtos.ProjectSearchResponse{
		Projects: helpers.MapProjectsToProjectListItems(projects),
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}, nil
}

func (s *ProjectService) GetProjectByID(c context.Context, id uint) (*dtos.ProjectDetail, error) {
	project, err := s.projectRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.ErrProjectNotFound
		}
		return nil, appErrors.ErrInternalServerError
	}

	return helpers.MapProjectToProjectDetail(project), nil
}

func (s *ProjectService) CreateProject(c context.Context, req dtos.CreateOrUpdateProjectRequest) error {
	var startDate, endDate *time.Time
	if req.StartDate != nil {
		startDate = &req.StartDate.Time
	}
	if req.EndDate != nil {
		endDate = &req.EndDate.Time
	}

	project := &models.Project{
		Name:         strings.TrimSpace(req.Name),
		Abbreviation: strings.TrimSpace(req.Abbreviation),
		StartDate:    startDate,
		EndDate:      endDate,
		LeaderID:     req.LeaderID,
		TeamID:       req.TeamID,
	}

	if err := s.projectRepository.Create(s.db.WithContext(c), project, req.MemberIDs); err != nil {
		if appErrors.IsDuplicatedEntryError(err) {
			return appErrors.ErrProjectAlreadyExists
		}
		return appErrors.ErrInternalServerError
	}
	return nil
}

func (s *ProjectService) UpdateProject(c context.Context, id uint, req dtos.CreateOrUpdateProjectRequest) error {
	_, err := s.projectRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrProjectNotFound
		}
		return appErrors.ErrInternalServerError
	}

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		startDate = &req.StartDate.Time
	}
	if req.EndDate != nil {
		endDate = &req.EndDate.Time
	}

	project := &models.Project{
		ID:           id,
		Name:         strings.TrimSpace(req.Name),
		Abbreviation: strings.TrimSpace(req.Abbreviation),
		StartDate:    startDate,
		EndDate:      endDate,
		LeaderID:     req.LeaderID,
		TeamID:       req.TeamID,
	}

	if err := s.projectRepository.Update(s.db.WithContext(c), project, req.MemberIDs); err != nil {
		return appErrors.ErrInternalServerError
	}
	return nil
}

func (s *ProjectService) DeleteProject(c context.Context, id uint) error {
	_, err := s.projectRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrProjectNotFound
		}
		return appErrors.ErrInternalServerError
	}

	if err := s.projectRepository.Delete(s.db.WithContext(c), id); err != nil {
		return appErrors.ErrInternalServerError
	}

	return nil
}
