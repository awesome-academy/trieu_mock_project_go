package services

import (
	"context"
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

	projectDtos := make([]dtos.ProjectSummary, 0, len(projects))
	for _, project := range projects {
		projectDtos = append(projectDtos, dtos.ProjectSummary{
			ID:           project.ID,
			Name:         project.Name,
			Abbreviation: project.Abbreviation,
			StartDate:    project.StartDate,
			EndDate:      project.EndDate,
		})
	}
	return projectDtos
}
