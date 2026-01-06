package repositories

import (
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type TeamsRepository struct {
}

func NewTeamsRepository() *TeamsRepository {
	return &TeamsRepository{}
}

func (r *TeamsRepository) ListTeams(db *gorm.DB, limit, offset int) ([]models.Team, error) {
	var teams []models.Team
	result := db.
		Preload("Leader").
		Preload("Members").
		Preload("Projects").
		Limit(limit).
		Offset(offset).
		Find(&teams)
	if result.Error != nil {
		return nil, result.Error
	}
	return teams, nil
}

func (r *TeamsRepository) CountTeams(db *gorm.DB) (int64, error) {
	var count int64
	result := db.Model(&models.Team{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r *TeamsRepository) FindByID(db *gorm.DB, id uint) (*models.Team, error) {
	var team models.Team
	result := db.
		Preload("Leader").
		Preload("Members").
		Preload("Projects").
		First(&team, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}

func (r *TeamsRepository) FindAllTeamsSummary(db *gorm.DB) ([]models.Team, error) {
	var teams []models.Team
	result := db.
		Select("id", "name").
		Find(&teams)
	if result.Error != nil {
		return nil, result.Error
	}
	return teams, nil
}

func (r *TeamsRepository) FindAllTeamsWithLeader(db *gorm.DB) ([]models.Team, error) {
	var teams []models.Team
	result := db.Preload("Leader").Find(&teams)
	if result.Error != nil {
		return nil, result.Error
	}
	return teams, nil
}

func (r *TeamsRepository) Create(db *gorm.DB, team *models.Team) error {
	return db.Create(team).Error
}

func (r *TeamsRepository) Update(db *gorm.DB, team *models.Team) error {
	return db.Model(&models.Team{}).
		Where("id = ?", team.ID).
		Updates(map[string]interface{}{
			"name":        team.Name,
			"description": team.Description,
			"leader_id":   team.LeaderID,
		}).Error
}

func (r *TeamsRepository) Delete(db *gorm.DB, id uint) error {
	return db.Delete(&models.Team{}, id).Error
}

func (r *TeamsRepository) ExistsByLeaderID(db *gorm.DB, leaderID uint) (bool, error) {
	var count int64
	result := db.Model(&models.Team{}).
		Where("leader_id = ?", leaderID).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
