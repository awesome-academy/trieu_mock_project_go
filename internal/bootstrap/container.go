package bootstrap

import (
	"trieu_mock_project_go/internal/config"
	"trieu_mock_project_go/internal/handlers"
	"trieu_mock_project_go/internal/middlewares"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

type AppContainer struct {
	// Middlewares
	AdminAuthMiddleware gin.HandlerFunc
	JWTAuthMiddleware   gin.HandlerFunc
	CSRFMiddleware      gin.HandlerFunc
	DataLoader          *middlewares.DataLoader

	// Services
	ValidationService *services.ValidationService
	AuthService       *services.AuthService
	UserService       *services.UserService
	TeamsService      *services.TeamsService
	PositionService   *services.PositionService
	ProjectService    *services.ProjectService
	SkillService      *services.SkillService

	// Handlers
	AuthHandler        *handlers.AuthHandler
	DashboardHandler   *handlers.DashboardHandler
	UserProfileHandler *handlers.UserProfileHandler
	TeamsHandler       *handlers.TeamsHandler
	// Admin Handlers
	AdminAuthHandler        *handlers.AdminAuthHandler
	AdminDashboardHandler   *handlers.AdminDashboardHandler
	AdminUserHandler        *handlers.AdminUserHandler
	AdminPositionHandler    *handlers.AdminPositionHandler
	AdminSkillHandler       *handlers.AdminSkillHandler
	AdminTeamHandler        *handlers.AdminTeamHandler
	AdminProjectHandler     *handlers.AdminProjectHandler
	AdminActivityLogHandler *handlers.AdminActivityLogHandler
}

func NewAppContainer() *AppContainer {
	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	teamsRepo := repositories.NewTeamsRepository()
	teamMemberRepo := repositories.NewTeamMemberRepository()
	positionRepo := repositories.NewPositionRepository()
	projectRepo := repositories.NewProjectRepository()
	projectMemberRepo := repositories.NewProjectMemberRepository()
	skillRepo := repositories.NewSkillRepository()
	activityLogRepo := repositories.NewActivityLogRepository()

	// Initialize services
	activityLogService := services.NewActivityLogService(config.DB, activityLogRepo)
	validationService := services.NewValidationService(config.DB, teamMemberRepo)
	authService := services.NewAuthService(config.DB, userRepo, activityLogService)
	userService := services.NewUserService(config.DB, userRepo, teamsRepo, projectRepo, projectMemberRepo, teamMemberRepo, activityLogService)
	teamsService := services.NewTeamsService(config.DB, teamsRepo, teamMemberRepo, userRepo, projectRepo, projectMemberRepo, activityLogService)
	positionService := services.NewPositionService(config.DB, positionRepo, activityLogService)
	projectService := services.NewProjectService(config.DB, projectRepo, validationService, activityLogService)
	skillService := services.NewSkillService(config.DB, skillRepo, activityLogService)

	return &AppContainer{
		// Middlewares
		JWTAuthMiddleware:   middlewares.JWTAuthMiddleware(),
		AdminAuthMiddleware: middlewares.AdminAuthMiddleware(),
		CSRFMiddleware:      middlewares.CSRFMiddleware(),
		DataLoader:          middlewares.NewDataLoader(teamsService, positionService, skillService),

		// Services
		ValidationService: validationService,
		AuthService:       authService,
		UserService:       userService,
		TeamsService:      teamsService,
		PositionService:   positionService,
		ProjectService:    projectService,
		SkillService:      skillService,

		// Handlers
		AuthHandler:        handlers.NewAuthHandler(authService),
		DashboardHandler:   handlers.NewDashboardHandler(),
		UserProfileHandler: handlers.NewUserProfileHandler(userService),
		TeamsHandler:       handlers.NewTeamsHandler(teamsService),
		// Admin Handlers
		AdminAuthHandler:        handlers.NewAdminAuthHandler(authService, activityLogService),
		AdminDashboardHandler:   handlers.NewAdminDashboardHandler(userService),
		AdminUserHandler:        handlers.NewAdminUserHandler(userService, teamsService, positionService, skillService),
		AdminPositionHandler:    handlers.NewAdminPositionHandler(positionService),
		AdminSkillHandler:       handlers.NewAdminSkillHandler(skillService),
		AdminTeamHandler:        handlers.NewAdminTeamHandler(teamsService, userService),
		AdminProjectHandler:     handlers.NewAdminProjectHandler(projectService, teamsService, userService),
		AdminActivityLogHandler: handlers.NewAdminActivityLogHandler(activityLogService),
	}
}
