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

type ProjectService struct {
	db                  *gorm.DB
	projectRepository   *repositories.ProjectRepository
	userRepository      *repositories.UserRepository
	validationService   *ValidationService
	activityLogService  *ActivityLogService
	notificationService *NotificationService
	emailService        *EmailService
}

func NewProjectService(db *gorm.DB, projectRepository *repositories.ProjectRepository, userRepository *repositories.UserRepository, validationService *ValidationService, activityLogService *ActivityLogService, notificationService *NotificationService, emailService *EmailService) *ProjectService {
	return &ProjectService{db: db, projectRepository: projectRepository, userRepository: userRepository, validationService: validationService, activityLogService: activityLogService, notificationService: notificationService, emailService: emailService}
}

func (s *ProjectService) GetAllProjectSummary(c context.Context) []dtos.ProjectSummary {
	projects, err := s.projectRepository.FindAllProjectSummary(s.db.WithContext(c))
	if err != nil {
		return []dtos.ProjectSummary{}
	}

	return helpers.MapProjectsToProjectSummaries(projects)
}

func (s *ProjectService) SearchProjects(c context.Context, teamID *uint, limit, offset int) (*dtos.ProjectSearchResponse, error) {
	projects, totalCount, err := s.projectRepository.SearchProjects(s.db.WithContext(c), teamID, limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	return &dtos.ProjectSearchResponse{
		Projects: helpers.MapProjectsToProjectListItems(projects),
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}, nil
}

func (s *ProjectService) GetProjectByID(c context.Context, id uint) (*dtos.ProjectDetail, error) {
	project, err := s.projectRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.ErrProjectNotFound
		}
		return nil, appErrors.ErrInternalServerError
	}

	return helpers.MapProjectToProjectDetail(project), nil
}

func (s *ProjectService) CreateProject(c context.Context, req dtos.CreateOrUpdateProjectRequest) error {
	if appErr := s.validationService.ValidateMembersInTeam(c, req.TeamID, req.MemberIDs); appErr != nil {
		return appErr
	}

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		startDate = &req.StartDate.Time
	}
	if req.EndDate != nil {
		endDate = &req.EndDate.Time
	}

	project := &models.Project{
		Name:         strings.TrimSpace(req.Name),
		Abbreviation: strings.TrimSpace(req.Abbreviation),
		StartDate:    startDate,
		EndDate:      endDate,
		LeaderID:     req.LeaderID,
		TeamID:       req.TeamID,
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.projectRepository.Create(tx, project, req.MemberIDs); err != nil {
			if appErrors.IsDuplicatedEntryError(err) {
				return appErrors.ErrProjectAlreadyExists
			}
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.CreateProject, project.ID, project.Name); err != nil {
			return err
		}

		if err := s.notificationService.NotifyProjectCreated(c, tx, project, req.MemberIDs); err != nil {
			return err
		}

		// Send email notifications to all members
		members, err := s.userRepository.FindByIDs(tx, req.MemberIDs)
		if err == nil {
			for _, m := range members {
				s.emailService.SendProjectJoinEmail(dtos.ProjectMembershipEmailDTO{
					To:          m.Email,
					UserName:    m.Name,
					ProjectName: project.Name,
				})
			}
		}

		return nil
	})
}

func (s *ProjectService) UpdateProject(c context.Context, id uint, req dtos.CreateOrUpdateProjectRequest) error {
	currentProject, err := s.projectRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrProjectNotFound
		}
		return appErrors.ErrInternalServerError
	}

	if appErr := s.validationService.ValidateMembersInTeam(c, req.TeamID, req.MemberIDs); appErr != nil {
		return appErr
	}

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		startDate = &req.StartDate.Time
	}
	if req.EndDate != nil {
		endDate = &req.EndDate.Time
	}

	projectToUpdate := &models.Project{
		ID:           id,
		Name:         strings.TrimSpace(req.Name),
		Abbreviation: strings.TrimSpace(req.Abbreviation),
		StartDate:    startDate,
		EndDate:      endDate,
		LeaderID:     req.LeaderID,
		TeamID:       req.TeamID,
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.projectRepository.Update(tx, projectToUpdate, req.MemberIDs); err != nil {
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.UpdateProject, projectToUpdate.ID, projectToUpdate.Name); err != nil {
			return err
		}

		currentMemberIDs := make([]uint, 0, len(currentProject.Members))
		for _, member := range currentProject.Members {
			currentMemberIDs = append(currentMemberIDs, member.ID)
		}
		if err := s.notificationService.NotifyProjectUpdated(c, tx, currentProject, currentMemberIDs, projectToUpdate, req.MemberIDs); err != nil {
			return err
		}

		// Send email notifications for added/removed members
		currentSet := utils.NewSet[uint]()
		for _, id := range currentMemberIDs {
			currentSet.Add(id)
		}
		updatedSet := utils.NewSet[uint]()
		for _, id := range req.MemberIDs {
			updatedSet.Add(id)
		}

		addedIDs := make([]uint, 0)
		for id := range updatedSet {
			if !currentSet.Has(id) {
				addedIDs = append(addedIDs, id)
			}
		}
		removedIDs := make([]uint, 0)
		for id := range currentSet {
			if !updatedSet.Has(id) {
				removedIDs = append(removedIDs, id)
			}
		}

		if len(addedIDs) > 0 {
			addedMembers, _ := s.userRepository.FindByIDs(tx, addedIDs)
			for _, m := range addedMembers {
				s.emailService.SendProjectJoinEmail(dtos.ProjectMembershipEmailDTO{
					To:          m.Email,
					UserName:    m.Name,
					ProjectName: projectToUpdate.Name,
				})
			}
		}
		if len(removedIDs) > 0 {
			removedMembers, _ := s.userRepository.FindByIDs(tx, removedIDs)
			for _, m := range removedMembers {
				s.emailService.SendProjectLeaveEmail(dtos.ProjectMembershipEmailDTO{
					To:          m.Email,
					UserName:    m.Name,
					ProjectName: projectToUpdate.Name,
				})
			}
		}

		return nil
	})
}

func (s *ProjectService) DeleteProject(c context.Context, id uint) error {
	project, err := s.projectRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrProjectNotFound
		}
		return appErrors.ErrInternalServerError
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.projectRepository.Delete(tx, id); err != nil {
			return appErrors.ErrInternalServerError
		}

		if err := s.activityLogService.LogActivityDb(c, tx, types.DeleteProject, project.ID, project.Name); err != nil {
			return err
		}

		if err := s.notificationService.NotifyProjectDeleted(c, tx, project.ID, project.Name); err != nil {
			return err
		}
		return nil
	})
}

func (s *ProjectService) ExportProjectsToCSV(c context.Context) ([][]string, error) {
	projects, err := s.projectRepository.FindAllProjectsWithMembers(s.db.WithContext(c))
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	data := [][]string{{"ID", "Name", "Abbreviation", "Start Date", "End Date", "LeaderId", "LeaderName", "TeamId", "TeamName", "MemberId", "MemberName"}}
	for _, p := range projects {
		startDate := ""
		if p.StartDate != nil {
			startDate = p.StartDate.Format("2006-01-02")
		}
		endDate := ""
		if p.EndDate != nil {
			endDate = p.EndDate.Format("2006-01-02")
		}
		projectBasicInfo := []string{
			fmt.Sprintf("%d", p.ID),
			p.Name,
			p.Abbreviation,
			startDate,
			endDate,
			fmt.Sprintf("%d", p.Leader.ID),
			p.Leader.Name,
			fmt.Sprintf("%d", p.Team.ID),
			p.Team.Name,
		}
		if len(p.Members) == 0 {
			data = append(data, append(projectBasicInfo, "", ""))
			continue
		} else {
			for _, member := range p.Members {
				data = append(data, append(projectBasicInfo, fmt.Sprintf("%d", member.ID), member.Name))
			}
		}
	}
	return data, nil
}

func (s *ProjectService) ImportProjectsFromCSV(c context.Context, data [][]string) error {
	if len(data) <= 1 {
		return appErrors.ErrNoCSVDataToImport
	}

	projectsMap := make(map[string]*dtos.ProjectImportData)
	projectNameSet := utils.NewSet[string]()
	teamIDSet := utils.NewSet[uint]()
	userIDSet := utils.NewSet[uint]()

	for i, row := range data {
		if i == 0 {
			continue
		}
		if len(row) < 10 {
			return fmt.Errorf("row %d: invalid number of columns", i+1)
		}

		name := strings.TrimSpace(row[0])
		abbreviation := strings.TrimSpace(row[1])
		startDateStr := strings.TrimSpace(row[2])
		endDateStr := strings.TrimSpace(row[3])
		leaderIDStr := strings.TrimSpace(row[4])
		teamIDStr := strings.TrimSpace(row[6])
		memberIDStr := strings.TrimSpace(row[8])

		if name == "" {
			return fmt.Errorf("row %d: project name is required", i+1)
		}

		rowNumber := i + 1

		p, ok := projectsMap[name]
		if !ok {
			var startDate, endDate *time.Time
			if startDateStr != "" {
				t, err := time.Parse("2006-01-02", startDateStr)
				if err == nil {
					startDate = &t
				}
			}
			if endDateStr != "" {
				t, err := time.Parse("2006-01-02", endDateStr)
				if err == nil {
					endDate = &t
				}
			}

			var leaderID uint
			if _, err := fmt.Sscanf(leaderIDStr, "%d", &leaderID); err != nil {
				return fmt.Errorf("row %d: invalid leader ID", rowNumber)
			}

			var teamID uint
			if _, err := fmt.Sscanf(teamIDStr, "%d", &teamID); err != nil {
				return fmt.Errorf("row %d: invalid team ID", rowNumber)
			}

			p = &dtos.ProjectImportData{
				Name:         name,
				Abbreviation: abbreviation,
				StartDate:    startDate,
				EndDate:      endDate,
				LeaderID:     leaderID,
				TeamID:       teamID,
				MemberIDs:    []uint{},
			}
			projectsMap[name] = p

			projectNameSet.Add(name)
			teamIDSet.Add(teamID)
			userIDSet.Add(leaderID)
		}

		if memberIDStr != "" {
			var memberID uint
			if _, err := fmt.Sscanf(memberIDStr, "%d", &memberID); err != nil {
				return fmt.Errorf("row %d: invalid member ID", rowNumber)
			}
			p.MemberIDs = append(p.MemberIDs, memberID)
			userIDSet.Add(memberID)
		}
	}

	for name, p := range projectsMap {
		if err := s.validationService.ValidateMembersInTeam(c, p.TeamID, p.MemberIDs); err != nil {
			return fmt.Errorf("Project '%s': %s", name, err.Error())
		}
	}

	if err := s.validationService.validateUserIDs(userIDSet.ToSlice()); err != nil {
		return err
	}

	if err := s.validationService.validateTeamIDs(teamIDSet.ToSlice()); err != nil {
		return err
	}

	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		activityLogs := make([]models.ActivityLog, 0, len(projectNameSet))

		for _, p := range projectsMap {
			project := &models.Project{
				Name:         p.Name,
				Abbreviation: p.Abbreviation,
				StartDate:    p.StartDate,
				EndDate:      p.EndDate,
				LeaderID:     p.LeaderID,
				TeamID:       p.TeamID,
			}

			if err := s.projectRepository.Create(tx, project, p.MemberIDs); err != nil {
				if appErrors.IsDuplicatedEntryError(err) {
					return appErrors.ErrProjectAlreadyExists
				}
				return appErrors.ErrInternalServerError
			}

			activityLog, err := s.activityLogService.createLogActivityModel(c, types.CreateProject, project.ID, project.Name)
			if err != nil {
				return err
			}
			activityLogs = append(activityLogs, *activityLog)

			// Send email notifications to all members
			members, errRepo := s.userRepository.FindByIDs(tx, p.MemberIDs)
			if errRepo == nil {
				for _, m := range members {
					s.emailService.SendProjectJoinEmail(dtos.ProjectMembershipEmailDTO{
						To:          m.Email,
						UserName:    m.Name,
						ProjectName: project.Name,
					})
				}
			}
		}

		if err := s.activityLogService.createInBatches(tx, activityLogs, 100); err != nil {
			return appErrors.ErrInternalServerError
		}

		return nil
	})
}
