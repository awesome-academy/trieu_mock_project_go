package repositories

import (
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type TeamMemberRepository struct {
}

func NewTeamMemberRepository() *TeamMemberRepository {
	return &TeamMemberRepository{}
}

func (r *TeamMemberRepository) FindMembersByTeamID(db *gorm.DB, teamID uint, limit, offset int) ([]models.TeamMember, error) {
	var members []models.TeamMember
	result := db.
		Preload("User").
		Where("team_id = ?", teamID).
		Limit(limit).
		Offset(offset).
		Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (r *TeamMemberRepository) CountMembersByTeamID(db *gorm.DB, teamID uint) (int64, error) {
	var count int64
	result := db.Model(&models.TeamMember{}).
		Where("team_id = ?", teamID).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
