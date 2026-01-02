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

	// Services
	AuthService     *services.AuthService
	UserService     *services.UserService
	TeamsService    *services.TeamsService
	PositionService *services.PositionService
	ProjectService  *services.ProjectService
	SkillService    *services.SkillService

	// Handlers
	AuthHandler        *handlers.AuthHandler
	DashboardHandler   *handlers.DashboardHandler
	UserProfileHandler *handlers.UserProfileHandler
	TeamsHandler       *handlers.TeamsHandler
	// Admin Handlers
	AdminAuthHandler      *handlers.AdminAuthHandler
	AdminDashboardHandler *handlers.AdminDashboardHandler
	AdminUserHandler      *handlers.AdminUserHandler
	AdminPositionHandler  *handlers.AdminPositionHandler
	AdminSkillHandler     *handlers.AdminSkillHandler
	AdminTeamHandler      *handlers.AdminTeamHandler
}

func NewAppContainer() *AppContainer {
	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	teamsRepo := repositories.NewTeamsRepository()
	teamMemberRepo := repositories.NewTeamMemberRepository()
	positionRepo := repositories.NewPositionRepository()
	projectRepo := repositories.NewProjectRepository()
	skillRepo := repositories.NewSkillRepository()

	// Initialize services
	authService := services.NewAuthService(config.DB, userRepo)
	userService := services.NewUserService(config.DB, userRepo, teamsRepo)
	teamsService := services.NewTeamsService(config.DB, teamsRepo, teamMemberRepo, userRepo)
	positionService := services.NewPositionService(config.DB, positionRepo)
	projectService := services.NewProjectService(config.DB, projectRepo)
	skillService := services.NewSkillService(config.DB, skillRepo)

	return &AppContainer{
		// Middlewares
		JWTAuthMiddleware:   middlewares.JWTAuthMiddleware(),
		AdminAuthMiddleware: middlewares.AdminAuthMiddleware(),
		CSRFMiddleware:      middlewares.CSRFMiddleware(),

		// Services
		AuthService:     authService,
		UserService:     userService,
		TeamsService:    teamsService,
		PositionService: positionService,
		ProjectService:  projectService,
		SkillService:    skillService,

		// Handlers
		AuthHandler:        handlers.NewAuthHandler(authService),
		DashboardHandler:   handlers.NewDashboardHandler(),
		UserProfileHandler: handlers.NewUserProfileHandler(userService),
		TeamsHandler:       handlers.NewTeamsHandler(teamsService),
		// Admin Handlers
		AdminAuthHandler:      handlers.NewAdminAuthHandler(authService),
		AdminDashboardHandler: handlers.NewAdminDashboardHandler(userService),
		AdminUserHandler:      handlers.NewAdminUserHandler(userService, teamsService, positionService, skillService),
		AdminPositionHandler:  handlers.NewAdminPositionHandler(positionService),
		AdminSkillHandler:     handlers.NewAdminSkillHandler(skillService),
		AdminTeamHandler:      handlers.NewAdminTeamHandler(teamsService, userService),
	}
}
