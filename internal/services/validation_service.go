package services

import (
	"context"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"

	"gorm.io/gorm"
)

type ValidationService struct {
	db                   *gorm.DB
	teamMemberRepository *repositories.TeamMemberRepository
}

func NewValidationService(db *gorm.DB, teamMemberRepository *repositories.TeamMemberRepository) *ValidationService {
	return &ValidationService{db: db, teamMemberRepository: teamMemberRepository}
}

func (s *ValidationService) ValidateMembersInTeam(c context.Context, teamID uint, memberIDs []uint) *appErrors.AppError {
	if len(memberIDs) == 0 {
		return nil
	}

	count, err := s.teamMemberRepository.CountActiveMembersInTeamByUserIDs(s.db.WithContext(c), teamID, memberIDs)
	if err != nil {
		return appErrors.ErrInternalServerError
	}

	if int(count) != len(memberIDs) {
		return appErrors.ErrUserNotInTeam
	}

	return nil
}
