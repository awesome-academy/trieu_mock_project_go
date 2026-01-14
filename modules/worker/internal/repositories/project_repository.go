package repositories

import (
	"time"
	"trieu_mock_project_go_worker/internal/models"

	"gorm.io/gorm"
)

type ProjectRepository struct{}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{}
}

func (r *ProjectRepository) FindProjectsNearDeadline(db *gorm.DB, days int) ([]models.Project, error) {
	var projects []models.Project
	now := time.Now()
	deadline := now.AddDate(0, 0, days)

	err := db.Preload("Members").
		Where("end_date IS NOT NULL AND end_date >= ? AND end_date <= ?",
			now.Format("2006-01-02"),
			deadline.Format("2006-01-02")).
		Find(&projects).Error
	return projects, err
}
