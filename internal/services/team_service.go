package services

import (
	"context"
	"fmt"
	"strings"
	"time"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/internal/types"
	"trieu_mock_project_go/internal/utils"
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
	notificationService     *NotificationService
	emailService            *EmailService
}

func NewTeamsService(db *gorm.DB,
	teamRepository *repositories.TeamsRepository,
	teamMemberRepository *repositories.TeamMemberRepository,
	userRepository *repositories.UserRepository,
	projectRepository *repositories.ProjectRepository,
	projectMemberRepository *repositories.ProjectMemberRepository,
	activityLogService *ActivityLogService,
	notificationService *NotificationService,
	emailService *EmailService) *TeamsService {
	return &TeamsService{
		db:                      db,
		teamRepository:          teamRepository,
		teamMemberRepository:    teamMemberRepository,
		userRepository:          userRepository,
		projectRepository:       projectRepository,
		projectMemberRepository: projectMemberRepository,
		activityLogService:      activityLogService,
		notificationService:     notificationService,
		emailService:            emailService,
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

		if err := s.activityLogService.LogActivityDb(c, tx, types.CreateTeam, team.ID, team.Name); err != nil {
			return err
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.JoinTeam, leader.ID, leader.Email, team.ID); err != nil {
			return err
		}

		if err := s.notificationService.NotifyTeamCreated(c, tx, team, leader.Name); err != nil {
			return err
		}

		// Send email notification to leader
		s.emailService.SendTeamJoinEmail(dtos.TeamMembershipEmailDTO{
			To:       leader.Email,
			UserName: leader.Name,
			TeamName: team.Name,
		})

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

	isNameChanged := team.Name != req.Name
	isDescriptionChanged := team.Description != req.Description
	isLeaderChanged := team.LeaderID != req.LeaderID
	if !isNameChanged && !isDescriptionChanged && !isLeaderChanged {
		return appErrors.ErrNoChangesDetected
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
					return err
				}

				// Send email notification to new leader
				s.emailService.SendTeamJoinEmail(dtos.TeamMembershipEmailDTO{
					To:       newLeader.Email,
					UserName: newLeader.Name,
					TeamName: team.Name,
				})
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
			return err
		}

		if err := s.notificationService.NotifyTeamUpdated(c, tx, team.ID, team.Name, isNameChanged || isDescriptionChanged, isLeaderChanged); err != nil {
			return err
		}
		return nil
	})
}

func (s *TeamsService) DeleteTeam(c context.Context, id uint) error {
	team, err := s.teamRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrTeamNotFound
		}
		return appErrors.ErrInternalServerError
	}
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Set current_team_id = null for all users in this team
		if err := s.userRepository.UpdateUsersCurrentTeamToNullByTeamID(tx, id); err != nil {
			return appErrors.ErrInternalServerError
		}

		if err := s.teamRepository.Delete(tx, id); err != nil {
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.DeleteTeam, team.ID, team.Name); err != nil {
			return err
		}

		if err := s.notificationService.NotifyTeamDeleted(c, tx, team.ID, team.Name); err != nil {
			return err
		}
		return nil
	})
}

func (s *TeamsService) AddMemberToTeam(c context.Context, teamID uint, userID uint) error {
	team, err := s.teamRepository.FindByID(s.db.WithContext(c), teamID)
	if err != nil {
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
				return err
			}

			if err := s.notificationService.NotifyTeamMemberRemoved(c, tx, activeTeamMember.TeamID, user.ID, user.Name); err != nil {
				return err
			}

			// Send leave email for old team
			s.emailService.SendTeamLeaveEmail(dtos.TeamMembershipEmailDTO{
				To:       user.Email,
				UserName: user.Name,
				TeamName: activeTeam.Name,
			})
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
			return err
		}

		if err := s.notificationService.NotifyTeamMemberAdded(c, tx, teamID, user.ID, user.Name); err != nil {
			return err
		}

		// Send email notification
		s.emailService.SendTeamJoinEmail(dtos.TeamMembershipEmailDTO{
			To:       user.Email,
			UserName: user.Name,
			TeamName: team.Name,
		})

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
			return err
		}

		if err := s.notificationService.NotifyTeamMemberRemoved(c, tx, teamID, user.ID, user.Name); err != nil {
			return err
		}

		// Send email notification
		s.emailService.SendTeamLeaveEmail(dtos.TeamMembershipEmailDTO{
			To:       user.Email,
			UserName: user.Name,
			TeamName: team.Name,
		})

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

func (s *TeamsService) ImportTeamsFromCSV(c context.Context, data [][]string) error {
	if len(data) <= 1 {
		return appErrors.ErrNoCSVDataToImport
	}

	var teamsToImport []dtos.TeamImportData
	uniqueTeamNames := utils.NewSet[string]()
	leaderIDSet := utils.NewSet[uint]()

	for i, row := range data {
		if i == 0 {
			continue
		}
		rowNumber := i + 1
		if len(row) < 3 {
			return fmt.Errorf("row %d: invalid number of columns", rowNumber)
		}

		name := strings.TrimSpace(row[0])
		description := strings.TrimSpace(row[1])
		leaderIDStr := strings.TrimSpace(row[2])

		if name == "" || leaderIDStr == "" {
			return fmt.Errorf("row %d: name and leader ID are required", rowNumber)
		}

		var leaderID uint
		if _, err := fmt.Sscanf(leaderIDStr, "%d", &leaderID); err != nil {
			return fmt.Errorf("row %d: invalid leader ID", rowNumber)
		}

		if uniqueTeamNames.Has(name) {
			return fmt.Errorf("duplicate team name '%s' in CSV", name)
		}

		uniqueTeamNames.Add(name)
		var descPtr *string
		if description != "" {
			descPtr = &description
		}
		teamsToImport = append(teamsToImport, dtos.TeamImportData{
			Name:        name,
			Description: descPtr,
			LeaderID:    leaderID,
		})
		if !leaderIDSet.Has(leaderID) {
			leaderIDSet.Add(leaderID)
		}
	}

	leaders, err := s.userRepository.FindByIDs(s.db.WithContext(c), leaderIDSet.ToSlice())
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	if len(leaders) != leaderIDSet.Size() {
		return appErrors.ErrUserNotFound
	}
	leaderMap := make(map[uint]models.User)
	for _, leader := range leaders {
		leaderMap[leader.ID] = leader
	}

	activityLogs := make([]models.ActivityLog, 0, len(teamsToImport)*2)
	teamMembers := make([]models.TeamMember, 0, len(teamsToImport))
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		for _, t := range teamsToImport {
			newTeam := models.Team{
				Name:        t.Name,
				Description: t.Description,
				LeaderID:    t.LeaderID,
			}
			if err := s.teamRepository.Create(tx, &newTeam); err != nil {
				if appErrors.IsDuplicatedEntryError(err) {
					return appErrors.ErrTeamAlreadyExists
				}
				return appErrors.ErrInternalServerError
			}

			leader := leaderMap[newTeam.LeaderID]
			leader.CurrentTeamID = &newTeam.ID
			if err := s.userRepository.UpdateUser(tx, &leader); err != nil {
				return appErrors.ErrInternalServerError
			}

			teamMember := models.TeamMember{
				UserID:   newTeam.LeaderID,
				TeamID:   newTeam.ID,
				JoinedAt: time.Now(),
			}
			teamMembers = append(teamMembers, teamMember)

			joinLog, err := s.activityLogService.createLogActivityModel(c, types.JoinTeam, newTeam.LeaderID, leader.Email, newTeam.ID)
			if err != nil {
				return err
			}
			activityLogs = append(activityLogs, *joinLog)

			// Send email notification to leader
			s.emailService.SendTeamJoinEmail(dtos.TeamMembershipEmailDTO{
				To:       leader.Email,
				UserName: leader.Name,
				TeamName: newTeam.Name,
			})
		}

		if err := s.teamMemberRepository.CreateInBatches(tx, teamMembers, 100); err != nil {
			if appErrors.IsDuplicatedEntryError(err) {
				return appErrors.ErrTeamLeaderAlreadyInAnotherTeam
			}
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.createInBatches(tx, activityLogs, 100); err != nil {
			return appErrors.ErrInternalServerError
		}

		return nil
	})
}
