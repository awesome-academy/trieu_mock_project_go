package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/internal/utils"
	"trieu_mock_project_go/internal/websocket"
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type NotificationService struct {
	db                     *gorm.DB
	notificationRepository *repositories.NotificationRepository
	userRepository         *repositories.UserRepository
	teamMemberRepository   *repositories.TeamMemberRepository
	projectRepository      *repositories.ProjectRepository
	redisService           *RedisService
	hub                    *websocket.Hub
}

func NewNotificationService(db *gorm.DB, notificationRepository *repositories.NotificationRepository,
	userRepository *repositories.UserRepository,
	teamMemberRepository *repositories.TeamMemberRepository,
	projectRepository *repositories.ProjectRepository,
	redisService *RedisService,
	hub *websocket.Hub) *NotificationService {
	return &NotificationService{
		db:                     db,
		notificationRepository: notificationRepository,
		userRepository:         userRepository,
		teamMemberRepository:   teamMemberRepository,
		projectRepository:      projectRepository,
		redisService:           redisService,
		hub:                    hub,
	}
}

func (s *NotificationService) GetUserNotifications(c context.Context, userID uint, limit, offset int) (*dtos.NotificationListResponse, *appErrors.AppError) {
	notifications, total, err := s.notificationRepository.FindByUserID(s.db.WithContext(c), userID, limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	return &dtos.NotificationListResponse{
		Notifications: helpers.MapNotificationsToNotificationResponses(notifications),
		Total:         total,
	}, nil
}

func (s *NotificationService) GetUnreadCount(c context.Context, userID uint) (int64, *appErrors.AppError) {
	count, err := s.notificationRepository.CountUnreadByUserID(s.db.WithContext(c), userID)
	if err != nil {
		return 0, appErrors.ErrInternalServerError
	}
	return count, nil
}

func (s *NotificationService) MarkAsRead(c context.Context, userID uint, notificationID uint) *appErrors.AppError {
	if err := s.notificationRepository.UpdateNotificationAsRead(s.db.WithContext(c), userID, notificationID); err != nil {
		return appErrors.ErrInternalServerError
	}
	return nil
}

func (s *NotificationService) MarkAllAsRead(c context.Context, userID uint) *appErrors.AppError {
	if err := s.notificationRepository.UpdateAllNotificationsAsReadByUserID(s.db.WithContext(c), userID); err != nil {
		return appErrors.ErrInternalServerError
	}
	return nil
}

func (s *NotificationService) DeleteNotification(c context.Context, userID uint, notificationID uint) *appErrors.AppError {
	if err := s.notificationRepository.DeleteByUserIDAndNotificationID(s.db.WithContext(c), userID, notificationID); err != nil {
		return appErrors.ErrInternalServerError
	}
	return nil
}

func (s *NotificationService) CreateNotificationDb(c context.Context, tx *gorm.DB, userID uint, title, content string) *appErrors.AppError {
	notification := &models.Notification{
		UserID:  userID,
		Title:   title,
		Content: content,
		IsRead:  false,
	}
	if err := s.notificationRepository.Create(tx, notification); err != nil {
		return appErrors.ErrInternalServerError
	}

	s.pushNotification(notification)

	return nil
}

func (s *NotificationService) NotifyTeamCreated(c context.Context, tx *gorm.DB, team *models.Team, leaderName string) *appErrors.AppError {
	title := "New Team Created"
	content := fmt.Sprintf("Team '%s' has been created. Manager: %s. Created at: %s",
		team.Name, leaderName, team.CreatedAt)

	// Notify leader
	return s.CreateNotificationDb(c, tx, team.LeaderID, title, content)
}

func (s *NotificationService) NotifyTeamUpdated(c context.Context, tx *gorm.DB, teamID uint, teamName string, isTeamInfoChanged, isLeaderChanged bool) *appErrors.AppError {
	title := "Team Updated"
	content := fmt.Sprintf("Team '%s' has been updated. Changes: ", teamName)
	if isTeamInfoChanged {
		content += "Team information changed. "
	}
	if isLeaderChanged {
		content += "Team leader changed."
	}

	memberIDs, err := s.teamMemberRepository.FindAllActiveMemberIDsByTeamID(tx, teamID)
	if err != nil {
		return appErrors.ErrInternalServerError
	}

	notifications := make([]models.Notification, 0, len(memberIDs))
	for _, userID := range memberIDs {
		notifications = append(notifications, models.Notification{
			UserID:  userID,
			Title:   title,
			Content: content,
			IsRead:  false,
		})
	}

	if len(notifications) > 0 {
		if err := s.notificationRepository.CreateInBatches(tx, notifications, 100); err != nil {
			return appErrors.ErrInternalServerError
		}
		s.pushNotifications(notifications)
	}
	return nil
}

func (s *NotificationService) NotifyTeamDeleted(c context.Context, tx *gorm.DB, teamID uint, teamName string) *appErrors.AppError {
	title := "Team Deleted"
	content := fmt.Sprintf("Team '%s' has been deleted.", teamName)

	// Get all members before deletion
	memberIDs, err := s.teamMemberRepository.FindAllActiveMemberIDsByTeamID(tx, teamID)
	if err != nil {
		return appErrors.ErrInternalServerError
	}

	notifications := make([]models.Notification, 0, len(memberIDs))
	for _, userID := range memberIDs {
		notifications = append(notifications, models.Notification{
			UserID:  userID,
			Title:   title,
			Content: content,
			IsRead:  false,
		})
	}

	if len(notifications) > 0 {
		if err := s.notificationRepository.CreateInBatches(tx, notifications, 100); err != nil {
			return appErrors.ErrInternalServerError
		}
		s.pushNotifications(notifications)
	}
	return nil
}

func (s *NotificationService) NotifyTeamMemberAdded(c context.Context, tx *gorm.DB, teamID uint, joinedUserID uint, joinedUserName string) *appErrors.AppError {
	// Notify the joined user
	title := "Added to Team"
	content := fmt.Sprintf("You have been added to the team '%s'.", joinedUserName)
	if err := s.CreateNotificationDb(c, tx, joinedUserID, title, content); err != nil {
		return err
	}

	// Notify other team members
	memberIDs, err := s.teamMemberRepository.FindAllActiveMemberIDsByTeamID(tx, teamID)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	notifications := make([]models.Notification, 0, len(memberIDs))
	for _, userID := range memberIDs {
		if userID == joinedUserID {
			continue
		}
		notifications = append(notifications, models.Notification{
			UserID:  userID,
			Title:   "New Team Member",
			Content: fmt.Sprintf("User '%s' has joined your team.", joinedUserName),
			IsRead:  false,
		})
	}

	if len(notifications) > 0 {
		if err := s.notificationRepository.CreateInBatches(tx, notifications, 100); err != nil {
			return appErrors.ErrInternalServerError
		}
		s.pushNotifications(notifications)
	}
	return nil
}

func (s *NotificationService) NotifyTeamMemberRemoved(c context.Context, tx *gorm.DB, teamID uint, removedUserID uint, removedUserName string) *appErrors.AppError {
	// Notify the removed user
	title := "Removed from Team"
	content := fmt.Sprintf("You have been removed from the team '%s'.", removedUserName)
	if err := s.CreateNotificationDb(c, tx, removedUserID, title, content); err != nil {
		return err
	}

	// Notify other team members
	memberIDs, err := s.teamMemberRepository.FindAllActiveMemberIDsByTeamID(tx, teamID)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	notifications := make([]models.Notification, 0, len(memberIDs))
	for _, userID := range memberIDs {
		if userID == removedUserID {
			continue
		}
		notifications = append(notifications, models.Notification{
			UserID:  userID,
			Title:   "Team Member Removed",
			Content: fmt.Sprintf("User '%s' has been removed from your team.", removedUserName),
			IsRead:  false,
		})
	}

	if len(notifications) > 0 {
		if err := s.notificationRepository.CreateInBatches(tx, notifications, 100); err != nil {
			return appErrors.ErrInternalServerError
		}
		s.pushNotifications(notifications)
	}
	return nil
}

func (s *NotificationService) NotifyProjectCreated(c context.Context, tx *gorm.DB, project *models.Project, memberIDs []uint) *appErrors.AppError {
	title := "New Project Created"

	leader, err := s.userRepository.FindByID(tx, project.LeaderID)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	content := fmt.Sprintf("Project '%s' has been created. Manager: %s. Created at: %s",
		project.Name, leader.Name, project.CreatedAt)

	notifications := make([]models.Notification, 0, len(memberIDs))
	for _, userID := range memberIDs {
		notifications = append(notifications, models.Notification{
			UserID:  userID,
			Title:   title,
			Content: content,
			IsRead:  false,
		})
	}

	if len(notifications) > 0 {
		if err := s.notificationRepository.CreateInBatches(tx, notifications, 100); err != nil {
			return appErrors.ErrInternalServerError
		}
		s.pushNotifications(notifications)
	}
	return nil
}

func (s *NotificationService) NotifyProjectUpdated(c context.Context, tx *gorm.DB,
	currentProject *models.Project, currentMemberIDs []uint,
	updatedProject *models.Project, updatedMemberIDs []uint) *appErrors.AppError {
	currentSet := utils.NewSet[uint]()
	for _, id := range currentMemberIDs {
		currentSet.Add(id)
	}

	updatedSet := utils.NewSet[uint]()
	for _, id := range updatedMemberIDs {
		updatedSet.Add(id)
	}

	removedIDs := make([]uint, 0)
	for id := range currentSet {
		if !updatedSet.Has(id) {
			removedIDs = append(removedIDs, id)
		}
	}

	addedIDs := make([]uint, 0)
	for id := range updatedSet {
		if !currentSet.Has(id) {
			addedIDs = append(addedIDs, id)
		}
	}

	stayedIDs := make([]uint, 0)
	for id := range updatedSet {
		if currentSet.Has(id) {
			stayedIDs = append(stayedIDs, id)
		}
	}

	notifications := make([]models.Notification, 0)

	// 1. Notify removed members
	if len(removedIDs) > 0 {
		title := "Removed from Project"
		content := fmt.Sprintf("You have been removed from the project '%s'.", currentProject.Name)
		for _, userID := range removedIDs {
			notifications = append(notifications, models.Notification{
				UserID:  userID,
				Title:   title,
				Content: content,
				IsRead:  false,
			})
		}
	}

	// 2. Notify added members
	if len(addedIDs) > 0 {
		title := "Added to Project"
		content := fmt.Sprintf("You have been added to the project '%s'.", updatedProject.Name)
		for _, userID := range addedIDs {
			notifications = append(notifications, models.Notification{
				UserID:  userID,
				Title:   title,
				Content: content,
				IsRead:  false,
			})
		}
	}

	// 3. Notify staying members about changes
	if len(stayedIDs) > 0 {
		isStartDateChanged := s.isTimeChanged(currentProject.StartDate, updatedProject.StartDate)
		isEndDateChanged := s.isTimeChanged(currentProject.EndDate, updatedProject.EndDate)

		isInfoChanged := currentProject.Name != updatedProject.Name ||
			currentProject.Abbreviation != updatedProject.Abbreviation ||
			currentProject.LeaderID != updatedProject.LeaderID ||
			currentProject.TeamID != updatedProject.TeamID ||
			isStartDateChanged ||
			isEndDateChanged

		addedCount := len(addedIDs)
		removedCount := len(removedIDs)

		if isInfoChanged || addedCount > 0 || removedCount > 0 {
			title := "Project Information Updated"
			content := fmt.Sprintf("Project '%s' has been updated.", updatedProject.Name)
			if isInfoChanged {
				content += " Basic information has changed."
			}

			if addedCount > 0 {
				content += fmt.Sprintf(" %d member(s) have been added.", addedCount)
			}

			if removedCount > 0 {
				content += fmt.Sprintf(" %d member(s) have been removed.", removedCount)
			}

			for _, userID := range stayedIDs {
				notifications = append(notifications, models.Notification{
					UserID:  userID,
					Title:   title,
					Content: content,
					IsRead:  false,
				})
			}
		}
	}

	if len(notifications) > 0 {
		if err := s.notificationRepository.CreateInBatches(tx, notifications, 100); err != nil {
			return appErrors.ErrInternalServerError
		}
		s.pushNotifications(notifications)
	}

	return nil
}

func (s *NotificationService) NotifyProjectDeleted(c context.Context, tx *gorm.DB, projectID uint, projectName string) *appErrors.AppError {
	memberIDs, err := s.projectRepository.FindMemberIDsByProjectID(tx, projectID)
	if err != nil {
		return appErrors.ErrInternalServerError
	}

	title := "Project Deleted"
	content := fmt.Sprintf("Project '%s' has been deleted. Please contact project manager for next steps.", projectName)

	notifications := make([]models.Notification, 0, len(memberIDs))
	for _, userID := range memberIDs {
		notifications = append(notifications, models.Notification{
			UserID:  userID,
			Title:   title,
			Content: content,
			IsRead:  false,
		})
	}

	if len(notifications) > 0 {
		if err := s.notificationRepository.CreateInBatches(tx, notifications, 100); err != nil {
			return appErrors.ErrInternalServerError
		}
		s.pushNotifications(notifications)
	}
	return nil
}

func (s *NotificationService) isTimeChanged(currentTime, updatedTime *time.Time) bool {
	if currentTime == nil && updatedTime == nil {
		return false
	}
	if currentTime == nil || updatedTime == nil {
		return true
	}
	return !currentTime.Equal(*updatedTime)
}

func (s *NotificationService) pushNotification(notification *models.Notification) {
	msg := &websocket.NotificationMessage{
		UserID:  notification.UserID,
		Title:   notification.Title,
		Content: notification.Content,
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error marshaling notification: %v\n", err)
		return
	}

	if err := s.redisService.Publish(context.Background(), websocket.RedisNotificationChannel, msgData); err != nil {
		fmt.Printf("Error publishing notification to Redis: %v\n", err)
		// Fallback to local push if Redis fails
		s.hub.SendNotification(notification.UserID, msg)
	}
}

func (s *NotificationService) pushNotifications(notifications []models.Notification) {
	for _, n := range notifications {
		s.pushNotification(&n)
	}
}
