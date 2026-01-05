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
	userRepository       *repositories.UserRepository
	positionRepository   *repositories.PositionRepository
	skillRepository      *repositories.SkillRepository
	teamRepository       *repositories.TeamsRepository
}

func NewValidationService(
	db *gorm.DB,
	teamMemberRepository *repositories.TeamMemberRepository,
	userRepository *repositories.UserRepository,
	positionRepository *repositories.PositionRepository,
	skillRepository *repositories.SkillRepository,
	teamRepository *repositories.TeamsRepository,
) *ValidationService {
	return &ValidationService{
		db:                   db,
		teamMemberRepository: teamMemberRepository,
		userRepository:       userRepository,
		positionRepository:   positionRepository,
		skillRepository:      skillRepository,
		teamRepository:       teamRepository,
	}
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

func (s *ValidationService) validateUserIDs(userIDs []uint) *appErrors.AppError {
	return s.validateIDs(
		userIDs,
		s.userRepository.CountByIDs,
		appErrors.ErrUserNotFound,
	)
}

func (s *ValidationService) validatePositionIDs(positionIDs []uint) *appErrors.AppError {
	return s.validateIDs(
		positionIDs,
		s.positionRepository.CountByIDs,
		appErrors.ErrPositionNotFound,
	)
}

func (s *ValidationService) validateSkillIDs(skillIDs []uint) *appErrors.AppError {
	return s.validateIDs(
		skillIDs,
		s.skillRepository.CountByIDs,
		appErrors.ErrSkillNotFound,
	)
}

func (s *ValidationService) validateTeamIDs(teamIDs []uint) *appErrors.AppError {
	return s.validateIDs(
		teamIDs,
		s.teamRepository.CountByIDs,
		appErrors.ErrTeamNotFound,
	)
}

func (s *ValidationService) validateIDs(
	ids []uint,
	countFunc func(db *gorm.DB, ids []uint) (int64, error),
	notFoundErr *appErrors.AppError,
) *appErrors.AppError {
	if len(ids) == 0 {
		return nil
	}

	count, err := countFunc(s.db, ids)
	if err != nil {
		return appErrors.ErrInternalServerError
	}

	if int(count) != len(ids) {
		return notFoundErr
	}

	return nil
}
