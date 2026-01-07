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

type UserService struct {
	db                      *gorm.DB
	userRepository          *repositories.UserRepository
	teamRepository          *repositories.TeamsRepository
	projectRepository       *repositories.ProjectRepository
	projectMemberRepository *repositories.ProjectMemberRepository
	teamMemberRepository    *repositories.TeamMemberRepository
	activityLogService      *ActivityLogService
	validationService       *ValidationService
}

func NewUserService(db *gorm.DB,
	userRepository *repositories.UserRepository,
	teamRepository *repositories.TeamsRepository,
	projectRepository *repositories.ProjectRepository,
	projectMemberRepository *repositories.ProjectMemberRepository,
	teamMemberRepository *repositories.TeamMemberRepository,
	activityLogService *ActivityLogService,
	validationService *ValidationService) *UserService {
	return &UserService{
		db:                      db,
		userRepository:          userRepository,
		teamRepository:          teamRepository,
		projectRepository:       projectRepository,
		projectMemberRepository: projectMemberRepository,
		teamMemberRepository:    teamMemberRepository,
		activityLogService:      activityLogService,
		validationService:       validationService,
	}
}

func (s *UserService) GetUserProfile(c context.Context, id uint) (*dtos.UserProfile, error) {
	user, err := s.userRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	userProfile := helpers.MapUserToUserProfile(user)

	return userProfile, nil
}

func (s *UserService) SearchUsers(c context.Context, name *string, teamId *uint, limit, offset int) (*dtos.UserSearchResponse, error) {
	users, totalCount, err := s.userRepository.SearchUsers(s.db.WithContext(c), name, teamId, limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	response := &dtos.UserSearchResponse{
		Users: helpers.MapUsersToUserDataForSearches(users),
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}

	return response, nil
}

func (s *UserService) CreateUser(c context.Context, req dtos.CreateOrUpdateUserRequest) error {
	existedUser, err := s.userRepository.FindByEmail(s.db.WithContext(c), req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return appErrors.ErrInternalServerError
	}
	if existedUser != nil {
		return appErrors.ErrEmailAlreadyExists
	}

	var birthday *time.Time
	if req.Birthday != nil && !req.Birthday.Time.IsZero() {
		birthday = &req.Birthday.Time
	}

	user := &models.User{
		Name:          req.Name,
		Email:         req.Email,
		Password:      utils.GenerateDefaultHashedPassword(),
		Birthday:      birthday,
		PositionID:    req.PositionID,
		Role:          "user", // Admin role is not allowed to be created here
		CurrentTeamID: req.TeamID,
	}

	err = s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		err := s.userRepository.CreateUser(tx, user)
		if err != nil {
			return err
		}

		userSkills := make([]models.UserSkill, 0, len(req.Skills))
		for _, sReq := range req.Skills {
			userSkills = append(userSkills, models.UserSkill{
				UserID:         user.ID,
				SkillID:        sReq.ID,
				Level:          sReq.Level,
				UsedYearNumber: sReq.UsedYearNumber,
			})
		}

		if err := s.userRepository.CreateUserSkills(tx, userSkills); err != nil {
			return err
		}

		if user.CurrentTeamID != nil {
			teamMember := &models.TeamMember{
				UserID:   user.ID,
				TeamID:   *user.CurrentTeamID,
				JoinedAt: time.Now(),
			}
			if err := s.teamMemberRepository.Create(tx, teamMember); err != nil {
				return err
			}

			if err := s.activityLogService.LogActivityDb(c, tx, types.JoinTeam, user.ID, user.Email, *user.CurrentTeamID); err != nil {
				return err
			}
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.CreateUser, user.ID); err != nil {
			return appErrors.ErrInternalServerError
		}

		return nil
	})

	if err != nil {
		return appErrors.ErrInternalServerError
	}

	return nil
}

func (s *UserService) UpdateUser(c context.Context, id uint, req dtos.CreateOrUpdateUserRequest) error {
	var birthday *time.Time
	if req.Birthday != nil && !req.Birthday.Time.IsZero() {
		birthday = &req.Birthday.Time
	}

	currentUser, appErr := s.userRepository.FindByID(s.db.WithContext(c), id)
	if appErr != nil {
		if appErr == gorm.ErrRecordNotFound {
			return appErrors.ErrUserNotFound
		}
		return appErrors.ErrInternalServerError
	}

	if currentUser.Email != req.Email {
		existedUser, err := s.userRepository.FindByEmail(s.db.WithContext(c), req.Email)
		if err != nil && err != gorm.ErrRecordNotFound {
			return appErrors.ErrInternalServerError
		}
		if existedUser != nil {
			return appErrors.ErrEmailAlreadyExists
		}
	}

	user := &models.User{
		ID:            id,
		Name:          req.Name,
		Email:         req.Email,
		Birthday:      birthday,
		PositionID:    req.PositionID,
		Role:          currentUser.Role,
		CurrentTeamID: req.TeamID,
	}

	userSkills := make([]models.UserSkill, 0, len(req.Skills))
	for _, sReq := range req.Skills {
		userSkills = append(userSkills, models.UserSkill{
			UserID:         id,
			SkillID:        sReq.ID,
			Level:          sReq.Level,
			UsedYearNumber: sReq.UsedYearNumber,
		})
	}

	appErr = s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Handle team change
		isTeamChanged := false
		if currentUser.CurrentTeamID == nil && req.TeamID != nil {
			isTeamChanged = true
		} else if currentUser.CurrentTeamID != nil && req.TeamID == nil {
			isTeamChanged = true
		} else if currentUser.CurrentTeamID != nil && req.TeamID != nil && *currentUser.CurrentTeamID != *req.TeamID {
			isTeamChanged = true
		}

		if isTeamChanged {
			// Handle leaving old team
			activeMember, err := s.teamMemberRepository.FindActiveMemberByUserID(tx, id)
			if err != nil {
				return appErrors.ErrInternalServerError
			}

			if activeMember != nil {
				// Check if user is leader of their current team
				isLeader, err := s.teamRepository.ExistByLeaderID(tx, id)
				if err != nil {
					return appErrors.ErrInternalServerError
				}
				if isLeader {
					return appErrors.ErrCannotRemoveOrMoveTeamLeader
				}

				// Check if user is member of any project in their current team
				isProjectMember, err := s.projectMemberRepository.ExistByMemberIdAndTeamId(tx, id, activeMember.TeamID)
				if err != nil {
					return appErrors.ErrInternalServerError
				}
				if isProjectMember {
					return appErrors.ErrCannotRemoveOrMoveProjectMember
				}

				// Set left_at for old team member record
				now := time.Now()
				activeMember.LeftAt = &now
				if err := s.teamMemberRepository.Update(tx, activeMember); err != nil {
					return appErrors.ErrInternalServerError
				}

				if err := s.activityLogService.LogActivityDb(c, tx, types.LeaveTeam, user.ID, user.Email, activeMember.TeamID); err != nil {
					return appErrors.ErrInternalServerError
				}
			}

			// Handle joining new team
			if req.TeamID != nil {
				newMember := &models.TeamMember{
					UserID:   id,
					TeamID:   *req.TeamID,
					JoinedAt: time.Now(),
				}
				if err := s.teamMemberRepository.Create(tx, newMember); err != nil {
					return appErrors.ErrInternalServerError
				}

				if err := s.activityLogService.LogActivityDb(c, tx, types.JoinTeam, user.ID, user.Email, *req.TeamID); err != nil {
					return appErrors.ErrInternalServerError
				}
			}
		}

		if err := s.userRepository.UpdateUser(tx, user); err != nil {
			return appErrors.ErrInternalServerError
		}
		if err := s.userRepository.UpdateUserSkills(tx, id, userSkills); err != nil {
			return appErrors.ErrInternalServerError
		}
		if err := s.activityLogService.LogActivityDb(c, tx, types.UpdateUser, user.ID); err != nil {
			return appErrors.ErrInternalServerError
		}

		return nil
	})
	return appErr
}

func (s *UserService) DeleteUser(c context.Context, id uint) error {
	user, err := s.userRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrUserNotFound
		}
		return appErrors.ErrInternalServerError
	}

	exist, err := s.teamRepository.ExistByLeaderID(s.db.WithContext(c), id)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	if exist {
		return appErrors.ErrCannotDeleteUserBeingTeamLeader
	}

	exist, err = s.projectRepository.ExistByLeaderId(s.db.WithContext(c), id)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	if exist {
		return appErrors.ErrCannotDeleteUserBeingProjectLeader
	}

	exist, err = s.projectMemberRepository.ExistByMemberId(s.db.WithContext(c), id)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	if exist {
		return appErrors.ErrCannotDeleteUserBeingProjectMember
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.userRepository.DeleteUser(tx, id); err != nil {
			return appErrors.ErrInternalServerError
		}
		if err := s.activityLogService.LogActivityDb(c, tx, types.DeleteUser, user.ID, user.Email); err != nil {
			return appErrors.ErrInternalServerError
		}
		return nil
	})
}

func (s *UserService) ExportUsersToCSV(c context.Context) ([][]string, error) {
	users, err := s.userRepository.FindAllUsersWithSkills(s.db.WithContext(c))
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	data := [][]string{{"ID", "Name", "Email", "Birthday", "PositionId", "PositionName", "TeamId", "TeamName", "SkillId", "SkillName", "SkillLevel", "SkillUsedYearNumber"}}
	for _, u := range users {
		birthday := ""
		if u.Birthday != nil {
			birthday = u.Birthday.Format("2006-01-02")
		}
		teamName := ""
		if u.CurrentTeam != nil {
			teamName = u.CurrentTeam.Name
		}
		for _, userSkill := range u.UserSkill {
			data = append(data, []string{
				fmt.Sprintf("%d", u.ID),
				u.Name,
				u.Email,
				birthday,
				fmt.Sprintf("%d", u.Position.ID),
				u.Position.Name,
				fmt.Sprintf("%d", *u.CurrentTeamID),
				teamName,
				fmt.Sprintf("%d", userSkill.Skill.ID),
				userSkill.Skill.Name,
				fmt.Sprintf("%d", userSkill.Level),
				fmt.Sprintf("%d", userSkill.UsedYearNumber),
			})
		}
	}
	return data, nil
}

func (s *UserService) ImportUsersFromCSV(c context.Context, data [][]string) error {
	if len(data) <= 1 {
		return appErrors.ErrNoCSVDataToImport
	}

	usersMap := make(map[string]*dtos.UserImportData)

	uniqueEmails := utils.NewSet[string]()
	uniquePositionIDs := utils.NewSet[uint]()
	uniqueSkillIDs := utils.NewSet[uint]()
	uniqueTeamIDs := utils.NewSet[uint]()

	for i, row := range data {
		if i == 0 {
			continue
		}
		if len(row) < 11 {
			return appErrors.ErrInvalidCSVFormat
		}

		rowNumber := i + 1

		name := strings.TrimSpace(row[0])
		email := strings.TrimSpace(row[1])
		birthdayStr := strings.TrimSpace(row[2])
		positionIDStr := strings.TrimSpace(row[3])
		teamIDStr := strings.TrimSpace(row[5])
		skillIDStr := strings.TrimSpace(row[7])
		skillLevelStr := strings.TrimSpace(row[9])
		skillUsedYearStr := strings.TrimSpace(row[10])

		if email == "" {
			return fmt.Errorf("row %d: email is required", rowNumber)
		}

		u, ok := usersMap[email]
		if !ok {
			var birthday *time.Time
			if birthdayStr != "" {
				t, err := time.Parse("2006-01-02", birthdayStr)
				if err != nil {
					return fmt.Errorf("row %d: invalid birthday format", rowNumber)
				}
				birthday = &t
			}

			var positionID uint
			if _, err := fmt.Sscanf(positionIDStr, "%d", &positionID); err != nil {
				return fmt.Errorf("row %d: invalid position ID", rowNumber)
			}
			uniquePositionIDs.Add(positionID)

			var teamID *uint
			if teamIDStr != "" && teamIDStr != "0" {
				var tid uint
				if _, err := fmt.Sscanf(teamIDStr, "%d", &tid); err == nil {
					teamID = &tid
				}
			}
			if teamID != nil {
				uniqueTeamIDs.Add(*teamID)
			}

			u = &dtos.UserImportData{
				Name:       name,
				Email:      email,
				Birthday:   birthday,
				PositionID: positionID,
				TeamID:     teamID,
				Skills:     make([]dtos.SkillImportData, 0),
			}
			uniqueEmails.Add(email)
			usersMap[email] = u
		}

		if skillIDStr != "" && skillIDStr != "0" {
			var skillID uint
			if _, err := fmt.Sscanf(skillIDStr, "%d", &skillID); err != nil {
				return fmt.Errorf("row %d: invalid skill ID", rowNumber)
			}
			var level int
			if _, err := fmt.Sscanf(skillLevelStr, "%d", &level); err != nil {
				return fmt.Errorf("row %d: invalid skill level", rowNumber)
			}
			var usedYear int
			if _, err := fmt.Sscanf(skillUsedYearStr, "%d", &usedYear); err != nil {
				return fmt.Errorf("row %d: invalid skill used year number", rowNumber)
			}

			if skillID > 0 {
				u.Skills = append(u.Skills, dtos.SkillImportData{
					ID:             skillID,
					Level:          level,
					UsedYearNumber: usedYear,
				})
				uniqueSkillIDs.Add(skillID)
			}
		}
	}

	if err := s.validationService.validatePositionIDs(uniquePositionIDs.ToSlice()); err != nil {
		return err
	}

	if err := s.validationService.validateSkillIDs(uniqueSkillIDs.ToSlice()); err != nil {
		return err
	}

	if err := s.validationService.validateTeamIDs(uniqueTeamIDs.ToSlice()); err != nil {
		return err
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		for _, u := range usersMap {
			user := &models.User{
				Name:          u.Name,
				Email:         u.Email,
				Password:      utils.GenerateDefaultHashedPassword(),
				Birthday:      u.Birthday,
				PositionID:    u.PositionID,
				Role:          "user", // Admin role is not allowed to be created here
				CurrentTeamID: u.TeamID,
			}

			if err := s.userRepository.CreateUser(tx, user); err != nil {
				if appErrors.IsDuplicatedEntryError(err) {
					return appErrors.ErrEmailAlreadyExists
				}
				return appErrors.ErrInternalServerError
			}

			userSkills := make([]models.UserSkill, 0, len(u.Skills))
			for i := range u.Skills {
				userSkills = append(userSkills, models.UserSkill{
					UserID:         user.ID,
					SkillID:        u.Skills[i].ID,
					Level:          u.Skills[i].Level,
					UsedYearNumber: u.Skills[i].UsedYearNumber,
				})
			}

			if len(u.Skills) > 0 {
				if err := s.userRepository.CreateUserSkills(tx, userSkills); err != nil {
					return err
				}
			}

			if user.CurrentTeamID != nil {
				teamMember := &models.TeamMember{
					UserID:   user.ID,
					TeamID:   *user.CurrentTeamID,
					JoinedAt: time.Now(),
				}
				if err := s.teamMemberRepository.Create(tx, teamMember); err != nil {
					return err
				}
				if err := s.activityLogService.LogActivityDb(c, tx, types.JoinTeam, user.ID, user.Email, *user.CurrentTeamID); err != nil {
					return err
				}
			}

			if err := s.activityLogService.LogActivityDb(c, tx, types.CreateUser, user.ID); err != nil {
				return err
			}
		}
		return nil
	})
}
