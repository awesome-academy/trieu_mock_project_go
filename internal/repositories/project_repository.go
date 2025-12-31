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

func (r *ProjectRepository) SearchProjects(db *gorm.DB, teamID *uint, limit, offset int) ([]models.Project, int64, error) {
	var projects []models.Project
	query := db.Model(&models.Project{})

	if teamID != nil && *teamID > 0 {
		query = query.Where("team_id = ?", *teamID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Leader").
		Preload("Team").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, count, nil
}

func (r *ProjectRepository) FindByID(db *gorm.DB, id uint) (*models.Project, error) {
	var project models.Project
	result := db.
		Preload("Leader").
		Preload("Team").
		Preload("Members").
		First(&project, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (r *ProjectRepository) Create(db *gorm.DB, project *models.Project, memberIDs []uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(project).Error; err != nil {
			return err
		}

		if len(memberIDs) > 0 {
			var members []models.User
			if err := tx.Where("id IN ?", memberIDs).Find(&members).Error; err != nil {
				return err
			}
			if err := tx.Model(project).Association("Members").Replace(members); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *ProjectRepository) Update(db *gorm.DB, project *models.Project, memberIDs []uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Project{}).
			Where("id = ?", project.ID).
			Updates(map[string]interface{}{
				"name":         project.Name,
				"abbreviation": project.Abbreviation,
				"start_date":   project.StartDate,
				"end_date":     project.EndDate,
				"leader_id":    project.LeaderID,
				"team_id":      project.TeamID,
			}).Error; err != nil {
			return err
		}

		var members []models.User
		if len(memberIDs) > 0 {
			if err := tx.Where("id IN ?", memberIDs).Find(&members).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(project).Association("Members").Replace(members); err != nil {
			return err
		}

		return nil
	})
}

func (r *ProjectRepository) Delete(db *gorm.DB, id uint) error {
	return db.Delete(&models.Project{}, id).Error
}
