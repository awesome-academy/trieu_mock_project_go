package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"trieu_mock_project_go/internal/config"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/handlers"
	"trieu_mock_project_go/internal/middlewares"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/internal/services"
	"trieu_mock_project_go/internal/websocket"

	"github.com/gin-gonic/gin"
)

type AppContainer struct {
	// WebSocket Hub
	Hub *websocket.Hub

	// Middlewares
	AdminAuthMiddleware gin.HandlerFunc
	JWTAuthMiddleware   gin.HandlerFunc
	JWTAuthWSMiddleware gin.HandlerFunc
	CSRFMiddleware      gin.HandlerFunc

	// Services
	ValidationService   *services.ValidationService
	AuthService         *services.AuthService
	UserService         *services.UserService
	TeamsService        *services.TeamsService
	PositionService     *services.PositionService
	ProjectService      *services.ProjectService
	SkillService        *services.SkillService
	NotificationService *services.NotificationService
	EmailService        *services.EmailService
	RabbitMQService     *services.RabbitMQService

	// Handlers
	AuthHandler         *handlers.AuthHandler
	DashboardHandler    *handlers.DashboardHandler
	UserProfileHandler  *handlers.UserProfileHandler
	TeamsHandler        *handlers.TeamsHandler
	NotificationHandler *handlers.NotificationHandler
	// Admin Handlers
	AdminAuthHandler        *handlers.AdminAuthHandler
	AdminDashboardHandler   *handlers.AdminDashboardHandler
	AdminUserHandler        *handlers.AdminUserHandler
	AdminPositionHandler    *handlers.AdminPositionHandler
	AdminSkillHandler       *handlers.AdminSkillHandler
	AdminTeamHandler        *handlers.AdminTeamHandler
	AdminProjectHandler     *handlers.AdminProjectHandler
	AdminActivityLogHandler *handlers.AdminActivityLogHandler
	AdminExportCsvHandler   *handlers.AdminExportCsvHandler
}

func NewAppContainer() *AppContainer {
	// Initialize WebSocket Hub
	hub := websocket.NewHub()

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	teamsRepo := repositories.NewTeamsRepository()
	teamMemberRepo := repositories.NewTeamMemberRepository()
	positionRepo := repositories.NewPositionRepository()
	projectRepo := repositories.NewProjectRepository()
	projectMemberRepo := repositories.NewProjectMemberRepository()
	skillRepo := repositories.NewSkillRepository()
	activityLogRepo := repositories.NewActivityLogRepository()
	notificationRepo := repositories.NewNotificationRepository()

	// Initialize services
	rabbitMQService := services.NewRabbitMQService()
	emailService := services.NewEmailService(rabbitMQService)
	redisService := services.NewRedisService()
	activityLogService := services.NewActivityLogService(config.DB, activityLogRepo)
	notificationService := services.NewNotificationService(config.DB, notificationRepo, userRepo, teamMemberRepo, projectRepo, redisService, hub)
	validationService := services.NewValidationService(config.DB, teamMemberRepo, userRepo, positionRepo, skillRepo, teamsRepo)
	authService := services.NewAuthService(config.DB, userRepo, activityLogService, redisService)
	userService := services.NewUserService(config.DB, userRepo, teamsRepo, projectRepo, projectMemberRepo, teamMemberRepo, activityLogService, validationService)
	teamsService := services.NewTeamsService(config.DB, teamsRepo, teamMemberRepo, userRepo, projectRepo, projectMemberRepo, activityLogService, notificationService, emailService)
	positionService := services.NewPositionService(config.DB, positionRepo, activityLogService)
	projectService := services.NewProjectService(config.DB, projectRepo, userRepo, validationService, activityLogService, notificationService, emailService)
	skillService := services.NewSkillService(config.DB, skillRepo, activityLogService)

	return &AppContainer{
		Hub: hub,
		// Middlewares
		JWTAuthMiddleware:   middlewares.JWTAuthMiddleware(authService),
		JWTAuthWSMiddleware: middlewares.JWTAuthWSMiddleware(authService),
		AdminAuthMiddleware: middlewares.AdminAuthMiddleware(),
		CSRFMiddleware:      middlewares.CSRFMiddleware(),

		// Services
		ValidationService:   validationService,
		AuthService:         authService,
		UserService:         userService,
		TeamsService:        teamsService,
		PositionService:     positionService,
		ProjectService:      projectService,
		SkillService:        skillService,
		NotificationService: notificationService,
		EmailService:        emailService,
		RabbitMQService:     rabbitMQService,

		// Handlers
		AuthHandler:         handlers.NewAuthHandler(authService),
		DashboardHandler:    handlers.NewDashboardHandler(),
		UserProfileHandler:  handlers.NewUserProfileHandler(userService),
		TeamsHandler:        handlers.NewTeamsHandler(teamsService),
		NotificationHandler: handlers.NewNotificationHandler(notificationService, hub),
		// Admin Handlers
		AdminAuthHandler:        handlers.NewAdminAuthHandler(authService, activityLogService),
		AdminDashboardHandler:   handlers.NewAdminDashboardHandler(userService),
		AdminUserHandler:        handlers.NewAdminUserHandler(userService, teamsService, positionService, skillService),
		AdminPositionHandler:    handlers.NewAdminPositionHandler(positionService),
		AdminSkillHandler:       handlers.NewAdminSkillHandler(skillService),
		AdminTeamHandler:        handlers.NewAdminTeamHandler(teamsService, userService),
		AdminProjectHandler:     handlers.NewAdminProjectHandler(projectService, teamsService, userService),
		AdminActivityLogHandler: handlers.NewAdminActivityLogHandler(activityLogService),
		AdminExportCsvHandler:   handlers.NewAdminExportCsvHandler(userService, positionService, projectService, skillService, teamsService),
	}
}

func (c *AppContainer) StartSubscriptionForNotifications() {
	go c.Hub.SubscribeToRedis(context.Background(), config.RedisClient, websocket.RedisNotificationChannel)
}

func (c *AppContainer) StartEmailWorker() error {
	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := range maxRetries {
		err := c.RabbitMQService.ConsumeEmailJobs(func(body []byte) error {
			var job dtos.EmailJobDTO

			if err := json.Unmarshal(body, &job); err != nil {
				return err
			}

			return c.EmailService.SendEmail(job.To, job.Subject, job.TemplateName, job.Data)
		})

		if err == nil {
			return nil
		}

		log.Printf("Failed to start email worker (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryInterval)
		time.Sleep(retryInterval)
	}

	return fmt.Errorf("failed to start email worker after %d retries", maxRetries)
}

func (c *AppContainer) InitializeApp() error {
	// Start Redis subscription for user notifications
	c.StartSubscriptionForNotifications()

	// Start RabbitMQ email worker
	if err := c.StartEmailWorker(); err != nil {
		return fmt.Errorf("failed to start email worker: %w", err)
	}

	return nil
}
