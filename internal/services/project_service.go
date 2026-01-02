package services

import (
	"context"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/repositories"

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
