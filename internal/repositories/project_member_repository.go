package repositories

import (
	"gorm.io/gorm"
)

type ProjectMemberRepository struct {
}

func NewProjectMemberRepository() *ProjectMemberRepository {
	return &ProjectMemberRepository{}
}

func (r *ProjectMemberRepository) ExistByMemberId(db *gorm.DB, memberId uint) (bool, error) {
	var count int64
	err := db.Table("project_members").
		Where("user_id = ?", memberId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProjectMemberRepository) ExistByMemberIdAndTeamId(db *gorm.DB, memberId uint, teamId uint) (bool, error) {
	var count int64
	err := db.Table("project_members").
		Joins("JOIN projects ON projects.id = project_members.project_id").
		Where("project_members.user_id = ? AND projects.team_id = ?", memberId, teamId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
