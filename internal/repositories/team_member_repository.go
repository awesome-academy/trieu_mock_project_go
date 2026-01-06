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

func (r *TeamMemberRepository) FindActiveMembersByTeamID(db *gorm.DB, teamID uint, limit, offset int) ([]models.TeamMember, error) {
	var members []models.TeamMember
	result := db.
		Preload("User").
		Where("team_id = ? AND left_at IS NULL", teamID).
		Limit(limit).
		Offset(offset).
		Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (r *TeamMemberRepository) FindAllActiveMembersByTeamID(db *gorm.DB, teamID uint) ([]models.TeamMember, error) {
	var members []models.TeamMember
	result := db.
		Preload("User").
		Where("team_id = ? AND left_at IS NULL", teamID).
		Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (r *TeamMemberRepository) CountActiveMembersByTeamID(db *gorm.DB, teamID uint) (int64, error) {
	var count int64
	result := db.Model(&models.TeamMember{}).
		Where("team_id = ? AND left_at IS NULL", teamID).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r *TeamMemberRepository) FindAllActiveMemberIDsByTeamID(db *gorm.DB, teamID uint) ([]uint, error) {
	var userIDs []uint
	result := db.Model(&models.TeamMember{}).
		Where("team_id = ? AND left_at IS NULL", teamID).
		Pluck("user_id", &userIDs)
	if result.Error != nil {
		return nil, result.Error
	}
	return userIDs, nil
}

func (r *TeamMemberRepository) FindTeamMembersByTeamID(db *gorm.DB, teamID uint, limit, offset int) ([]models.TeamMember, error) {
	var members []models.TeamMember
	result := db.
		Preload("User").
		Where("team_id = ?", teamID).
		Order("joined_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (r *TeamMemberRepository) CountTeamMembersByTeamID(db *gorm.DB, teamID uint) (int64, error) {
	var count int64
	result := db.Model(&models.TeamMember{}).
		Where("team_id = ?", teamID).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r *TeamMemberRepository) FindActiveMemberByUserID(db *gorm.DB, userID uint) (*models.TeamMember, error) {
	var member models.TeamMember
	result := db.Where("user_id = ? AND left_at IS NULL", userID).First(&member)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &member, nil
}

func (r *TeamMemberRepository) Create(db *gorm.DB, member *models.TeamMember) error {
	return db.Create(member).Error
}

func (r *TeamMemberRepository) CreateInBatches(db *gorm.DB, members []models.TeamMember, batchSize int) error {
	return db.CreateInBatches(members, batchSize).Error
}

func (r *TeamMemberRepository) Update(db *gorm.DB, member *models.TeamMember) error {
	return db.Model(&models.TeamMember{}).
		Where("id = ?", member.ID).
		Updates(map[string]interface{}{
			"team_id":   member.TeamID,
			"user_id":   member.UserID,
			"joined_at": member.JoinedAt,
			"left_at":   member.LeftAt,
		}).Error
}

func (r *TeamMemberRepository) CountActiveMembersInTeamByUserIDs(db *gorm.DB, teamID uint, userIDs []uint) (int64, error) {
	var count int64
	result := db.Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id IN ? AND left_at IS NULL", teamID, userIDs).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
