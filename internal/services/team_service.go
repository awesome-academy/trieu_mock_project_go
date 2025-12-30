package services

import (
	"context"
	"time"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/models"

	appErrors "trieu_mock_project_go/internal/errors"

	"gorm.io/gorm"
)

type TeamsService struct {
	db                   *gorm.DB
	teamRepository       *repositories.TeamsRepository
	teamMemberRepository *repositories.TeamMemberRepository
	userRepository       *repositories.UserRepository
}

func NewTeamsService(db *gorm.DB, teamRepository *repositories.TeamsRepository, teamMemberRepository *repositories.TeamMemberRepository, userRepository *repositories.UserRepository) *TeamsService {
	return &TeamsService{db: db, teamRepository: teamRepository, teamMemberRepository: teamMemberRepository, userRepository: userRepository}
}

func (s *TeamsService) ListTeams(c context.Context, limit, offset int) (*dtos.ListTeamsResponse, error) {
	teams, err := s.teamRepository.ListTeams(s.db.WithContext(c), limit, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.teamRepository.CountTeams(s.db.WithContext(c))
	if err != nil {
		return nil, err
	}

	teamDtos := helpers.MapTeamsToTeamDtos(teams)

	response := &dtos.ListTeamsResponse{
		Teams: teamDtos,
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}

	return response, nil
}

func (s *TeamsService) GetTeamDetails(c context.Context, id uint) (*dtos.Team, error) {
	team, err := s.teamRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		return nil, err
	}

	return helpers.MapTeamToTeamDto(team), nil
}

func (s *TeamsService) GetTeamMembers(c context.Context, teamID uint, limit, offset int) (*dtos.ListTeamMembersResponse, error) {
	members, err := s.teamMemberRepository.FindMembersByTeamID(s.db.WithContext(c), teamID, limit, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.teamMemberRepository.CountMembersByTeamID(s.db.WithContext(c), teamID)
	if err != nil {
		return nil, err
	}

	response := &dtos.ListTeamMembersResponse{
		Members: helpers.MapTeamMembersToTeamMemberSummaries(members),
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}

	return response, nil
}

func (s *TeamsService) GetAllTeamsSummary(c context.Context) []dtos.TeamSummary {
	teams, err := s.teamRepository.FindAllTeamsSummary(s.db.WithContext(c))
	if err != nil {
		return []dtos.TeamSummary{}
	}

	return helpers.MapTeamsToTeamSummaries(teams)
}

func (s *TeamsService) CreateTeam(c context.Context, req dtos.CreateOrUpdateTeamRequest) error {
	leader, err := s.userRepository.FindByID(s.db.WithContext(c), req.LeaderID)
	if err != nil {
		return err
	}

	team := &models.Team{
		Name:        req.Name,
		Description: req.Description,
		LeaderID:    req.LeaderID,
	}
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.teamRepository.Create(tx, team); err != nil {
			if appErrors.IsDuplicatedEntryError(err) {
				return appErrors.ErrTeamAlreadyExists
			}
			return err
		}
		leader.CurrentTeamID = &team.ID
		if err := s.userRepository.UpdateUser(tx, leader); err != nil {
			return err
		}
		newMember := &models.TeamMember{
			UserID:   req.LeaderID,
			TeamID:   team.ID,
			JoinedAt: time.Now(),
		}
		if err := s.teamMemberRepository.Create(tx, newMember); err != nil {
			// For team_members.ux_active_user_in_team unique constraint
			if appErrors.IsDuplicatedEntryError(err) {
				return appErrors.ErrTeamLeaderAlreadyInAnotherTeam
			}
			return err
		}
		return nil
	})
}

func (s *TeamsService) UpdateTeam(c context.Context, id uint, req dtos.CreateOrUpdateTeamRequest) error {
	team, err := s.teamRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		return err
	}

	team.Name = req.Name
	team.Description = req.Description

	if team.LeaderID != req.LeaderID {
		team.LeaderID = req.LeaderID
		return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
			if err := s.teamRepository.Update(tx, team); err != nil {
				if appErrors.IsDuplicatedEntryError(err) {
					return appErrors.ErrTeamAlreadyExists
				}
				return err
			}

			// Update leader's current_team_id
			newLeader, err := s.userRepository.FindByID(tx, req.LeaderID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return appErrors.ErrUserNotFound
				}
				return err
			}
			newLeader.CurrentTeamID = &team.ID
			if err := s.userRepository.UpdateUser(tx, newLeader); err != nil {
				return err
			}

			// Add new leader as team member
			newMember := &models.TeamMember{
				UserID:   req.LeaderID,
				TeamID:   team.ID,
				JoinedAt: time.Now(),
			}
			if err := s.teamMemberRepository.Create(tx, newMember); err != nil {
				// For team_members.ux_active_user_in_team unique constraint
				if appErrors.IsDuplicatedEntryError(err) {
					return appErrors.ErrTeamLeaderAlreadyInAnotherTeam
				}
				return err
			}

			return nil
		})
	} else {
		err = s.teamRepository.Update(s.db.WithContext(c), team)
		if appErrors.IsDuplicatedEntryError(err) {
			return appErrors.ErrTeamAlreadyExists
		}
		return err
	}
}

func (s *TeamsService) DeleteTeam(c context.Context, id uint) error {
	_, err := s.teamRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		return err
	}
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Set current_team_id = null for all users in this team
		if err = s.userRepository.UpdateUsersCurrentTeamToNullByTeamID(tx, id); err != nil {
			return err
		}

		return s.teamRepository.Delete(tx, id)
	})
}

func (s *TeamsService) AddMemberToTeam(c context.Context, teamID uint, userID uint) error {
	if _, err := s.teamRepository.FindByID(s.db.WithContext(c), teamID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrTeamNotFound
		}
		return err
	}
	user, err := s.userRepository.FindByID(s.db.WithContext(c), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrUserNotFound
		}
		return err
	}
	if (user.CurrentTeamID != nil) && (*user.CurrentTeamID == teamID) {
		return appErrors.ErrUserAlreadyInTeam
	}
	activeTeamMember, err := s.teamMemberRepository.FindActiveMemberByUserID(s.db.WithContext(c), userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	now := time.Now()
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// If user is in another team, set left_at for left team member record
		if activeTeamMember != nil {
			activeTeamMember.LeftAt = &now
			if err := s.teamMemberRepository.Update(tx, activeTeamMember); err != nil {
				return err
			}
		}
		// Add new team member record
		newMember := &models.TeamMember{
			UserID:   userID,
			TeamID:   teamID,
			JoinedAt: now,
		}
		if err := s.teamMemberRepository.Create(tx, newMember); err != nil {
			return err
		}

		// Update user's current_team_id
		user.CurrentTeamID = &teamID
		return s.userRepository.UpdateUser(tx, user)
	})
}

func (s *TeamsService) RemoveMemberFromTeam(c context.Context, teamID uint, userID uint) error {
	if _, err := s.teamRepository.FindByID(s.db.WithContext(c), teamID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrTeamNotFound
		}
		return err
	}
	user, err := s.userRepository.FindByID(s.db.WithContext(c), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrUserNotFound
		}
		return err
	}
	if (user.CurrentTeamID == nil) || (*user.CurrentTeamID != teamID) {
		return appErrors.ErrUserNotInTeam
	}
	teamMember, err := s.teamMemberRepository.FindActiveMemberByUserID(s.db.WithContext(c), userID)
	if err != nil {
		return err
	}
	if teamMember == nil || teamMember.TeamID != teamID {
		return appErrors.ErrUserNotInTeam
	}
	now := time.Now()
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Set left_at for team member record
		teamMember.LeftAt = &now
		if err := s.teamMemberRepository.Update(tx, teamMember); err != nil {
			return err
		}

		// Set user's current_team_id to null
		user.CurrentTeamID = nil
		return s.userRepository.UpdateUser(tx, user)
	})
}
