package repositories

import (
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type ProjectRepository struct {
}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{}
}

func (r *ProjectRepository) FindAllProjectSummary(db *gorm.DB) ([]models.Project, error) {
	var projects []models.Project
	result := db.
		Select("id", "name", "abbreviation", "start_date", "end_date").
		Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}
