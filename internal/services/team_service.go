package services

import (
	"context"
	"fmt"
	"time"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/internal/types"
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type TeamsService struct {
	db                      *gorm.DB
	teamRepository          *repositories.TeamsRepository
	teamMemberRepository    *repositories.TeamMemberRepository
	userRepository          *repositories.UserRepository
	projectRepository       *repositories.ProjectRepository
	projectMemberRepository *repositories.ProjectMemberRepository
	activityLogService      *ActivityLogService
}

func NewTeamsService(db *gorm.DB,
	teamRepository *repositories.TeamsRepository,
	teamMemberRepository *repositories.TeamMemberRepository,
	userRepository *repositories.UserRepository,
	projectRepository *repositories.ProjectRepository,
	projectMemberRepository *repositories.ProjectMemberRepository,
	activityLogService *ActivityLogService) *TeamsService {
	return &TeamsService{
		db:                      db,
		teamRepository:          teamRepository,
		teamMemberRepository:    teamMemberRepository,
		userRepository:          userRepository,
		projectRepository:       projectRepository,
		projectMemberRepository: projectMemberRepository,
		activityLogService:      activityLogService,
	}
}

func (s *TeamsService) ListTeams(c context.Context, limit, offset int) (*dtos.ListTeamsResponse, error) {
	teams, err := s.teamRepository.ListTeams(s.db.WithContext(c), limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	totalCount, err := s.teamRepository.CountTeams(s.db.WithContext(c))
	if err != nil {
		return nil, appErrors.ErrInternalServerError
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
		return nil, appErrors.ErrInternalServerError
	}

	return helpers.MapTeamToTeamDto(team), nil
}

func (s *TeamsService) GetTeamMembers(c context.Context, teamID uint, limit, offset int) (*dtos.ListTeamMembersResponse, error) {
	members, err := s.teamMemberRepository.FindActiveMembersByTeamID(s.db.WithContext(c), teamID, limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	totalCount, err := s.teamMemberRepository.CountActiveMembersByTeamID(s.db.WithContext(c), teamID)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
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

func (s *TeamsService) GetAllTeamMembers(c context.Context, teamID uint) ([]dtos.TeamMemberSummary, error) {
	members, err := s.teamMemberRepository.FindAllActiveMembersByTeamID(s.db.WithContext(c), teamID)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	return helpers.MapTeamMembersToTeamMemberSummaries(members), nil
}

func (s *TeamsService) GetTeamMemberHistory(c context.Context, teamID uint, limit, offset int) (*dtos.ListTeamMemberHistoryResponse, error) {
	members, err := s.teamMemberRepository.FindTeamMembersByTeamID(s.db.WithContext(c), teamID, limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	totalCount, err := s.teamMemberRepository.CountTeamMembersByTeamID(s.db.WithContext(c), teamID)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	return &dtos.ListTeamMemberHistoryResponse{
		History: helpers.MapTeamMembersToTeamMemberHistories(members),
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}, nil
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
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrUserNotFound
		}
		return appErrors.ErrInternalServerError
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
			return appErrors.ErrInternalServerError
		}
		leader.CurrentTeamID = &team.ID
		if err := s.userRepository.UpdateUser(tx, leader); err != nil {
			return appErrors.ErrInternalServerError
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
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.JoinTeam, leader.ID, leader.Email, team.ID); err != nil {
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.CreateTeam, team.ID, team.Name); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *TeamsService) UpdateTeam(c context.Context, id uint, req dtos.CreateOrUpdateTeamRequest) error {
	team, err := s.teamRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrTeamNotFound
		}
		return appErrors.ErrInternalServerError
	}

	team.Name = req.Name
	team.Description = req.Description

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if team.LeaderID != req.LeaderID {
			team.LeaderID = req.LeaderID
			if err := s.teamRepository.Update(tx, team); err != nil {
				if appErrors.IsDuplicatedEntryError(err) {
					return appErrors.ErrTeamAlreadyExists
				}
				return appErrors.ErrInternalServerError
			}

			// Update leader's current_team_id
			newLeader, err := s.userRepository.FindByID(tx, req.LeaderID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return appErrors.ErrUserNotFound
				}
				return appErrors.ErrInternalServerError
			}
			newLeader.CurrentTeamID = &team.ID
			if err := s.userRepository.UpdateUser(tx, newLeader); err != nil {
				return appErrors.ErrInternalServerError
			}

			activeTeamMember, err := s.teamMemberRepository.FindActiveMemberByUserID(tx, req.LeaderID)
			if err != nil && err != gorm.ErrRecordNotFound {
				return appErrors.ErrInternalServerError
			}
			if activeTeamMember != nil {
				if activeTeamMember.TeamID != team.ID {
					return appErrors.ErrTeamLeaderAlreadyInAnotherTeam
				}
			} else {
				// Add new leader as team member
				newMember := &models.TeamMember{
					UserID:   req.LeaderID,
					TeamID:   team.ID,
					JoinedAt: time.Now(),
				}
				if err := s.teamMemberRepository.Create(tx, newMember); err != nil {
					return appErrors.ErrInternalServerError
				}

				if err := s.activityLogService.LogActivityDb(c, tx, types.JoinTeam, newLeader.ID, newLeader.Email, team.ID); err != nil {
					return appErrors.ErrInternalServerError
				}
			}
		} else {
			if err = s.teamRepository.Update(tx, team); err != nil {
				if appErrors.IsDuplicatedEntryError(err) {
					return appErrors.ErrTeamAlreadyExists
				}
				return appErrors.ErrInternalServerError
			}
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.UpdateTeam, team.ID, team.Name); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *TeamsService) DeleteTeam(c context.Context, id uint) error {
	team_, err := s.teamRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrTeamNotFound
		}
		return appErrors.ErrInternalServerError
	}
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Set current_team_id = null for all users in this team
		if err := s.userRepository.UpdateUsersCurrentTeamToNullByTeamID(tx, id); err != nil {
			return err
		}

		if err := s.teamRepository.Delete(tx, id); err != nil {
			return err
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.DeleteTeam, team_.ID, team_.Name); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *TeamsService) AddMemberToTeam(c context.Context, teamID uint, userID uint) error {
	if _, err := s.teamRepository.FindByID(s.db.WithContext(c), teamID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrTeamNotFound
		}
		return appErrors.ErrInternalServerError
	}
	user, err := s.userRepository.FindByID(s.db.WithContext(c), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrUserNotFound
		}
		return appErrors.ErrInternalServerError
	}
	if (user.CurrentTeamID != nil) && (*user.CurrentTeamID == teamID) {
		return appErrors.ErrUserAlreadyInTeam
	}
	activeTeamMember, err := s.teamMemberRepository.FindActiveMemberByUserID(s.db.WithContext(c), userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return appErrors.ErrInternalServerError
	}
	now := time.Now()
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// If user is in another team, set left_at for left team member record
		if activeTeamMember != nil {
			activeTeam, err := s.teamRepository.FindByID(tx, activeTeamMember.TeamID)
			if err != nil {
				return appErrors.ErrInternalServerError
			}
			if activeTeam.LeaderID == userID {
				return appErrors.ErrCannotRemoveOrMoveTeamLeader
			}

			// Check if user is member of any project in current team
			isProjectMember, err := s.projectMemberRepository.ExistByMemberIdAndTeamId(tx, userID, activeTeamMember.TeamID)
			if err != nil {
				return appErrors.ErrInternalServerError
			}
			if isProjectMember {
				return appErrors.ErrCannotRemoveOrMoveProjectMember
			}

			activeTeamMember.LeftAt = &now
			if err := s.teamMemberRepository.Update(tx, activeTeamMember); err != nil {
				return appErrors.ErrInternalServerError
			}

			if err := s.activityLogService.LogActivityDb(c, tx, types.LeaveTeam, user.ID, user.Email, activeTeamMember.TeamID); err != nil {
				return appErrors.ErrInternalServerError
			}
		}
		// Add new team member record
		newMember := &models.TeamMember{
			UserID:   userID,
			TeamID:   teamID,
			JoinedAt: now,
		}
		if err := s.teamMemberRepository.Create(tx, newMember); err != nil {
			return appErrors.ErrInternalServerError
		}

		// Update user's current_team_id
		user.CurrentTeamID = &teamID
		if err := s.userRepository.UpdateUser(tx, user); err != nil {
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.JoinTeam, user.ID, user.Email, teamID); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *TeamsService) RemoveMemberFromTeam(c context.Context, teamID uint, userID uint) error {
	team, err := s.teamRepository.FindByID(s.db.WithContext(c), teamID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrTeamNotFound
		}
		return appErrors.ErrInternalServerError
	}
	if team.LeaderID == userID {
		return appErrors.ErrCannotRemoveOrMoveTeamLeader
	}
	user, err := s.userRepository.FindByID(s.db.WithContext(c), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrUserNotFound
		}
		return appErrors.ErrInternalServerError
	}
	if (user.CurrentTeamID == nil) || (*user.CurrentTeamID != teamID) {
		return appErrors.ErrUserNotInTeam
	}
	teamMember, err := s.teamMemberRepository.FindActiveMemberByUserID(s.db.WithContext(c), userID)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	if teamMember == nil || teamMember.TeamID != teamID {
		return appErrors.ErrUserNotInTeam
	}

	// Check if user is member of any project in this team
	isProjectMember, err := s.projectMemberRepository.ExistByMemberIdAndTeamId(s.db.WithContext(c), userID, teamID)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	if isProjectMember {
		return appErrors.ErrCannotRemoveOrMoveProjectMember
	}

	now := time.Now()
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Set left_at for team member record
		teamMember.LeftAt = &now
		if err := s.teamMemberRepository.Update(tx, teamMember); err != nil {
			return appErrors.ErrInternalServerError
		}

		// Set user's current_team_id to null
		user.CurrentTeamID = nil
		if err := s.userRepository.UpdateUser(tx, user); err != nil {
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.LeaveTeam, user.ID, user.Email, teamID); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *TeamsService) ExportTeamsToCSV(c context.Context) ([][]string, error) {
	teams, err := s.teamRepository.FindAllTeamsWithLeader(s.db.WithContext(c))
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	data := [][]string{{"ID", "Name", "Description", "LeaderId", "LeaderName"}}
	for _, t := range teams {
		description := ""
		if t.Description != nil {
			description = *t.Description
		}
		data = append(data, []string{
			fmt.Sprintf("%d", t.ID),
			t.Name,
			description,
			fmt.Sprintf("%d", t.LeaderID),
			t.Leader.Name,
		})
	}
	return data, nil
}
