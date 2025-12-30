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

	// Services
	AuthService *services.AuthService
	UserService *services.UserService
	TeamsService *services.TeamsService

	// Handlers
	AuthHandler        *handlers.AuthHandler
	DashboardHandler   *handlers.DashboardHandler
	UserProfileHandler *handlers.UserProfileHandler
	TeamsHandler       *handlers.TeamsHandler
	// Admin Handlers
	AdminAuthHandler      *handlers.AdminAuthHandler
	AdminDashboardHandler *handlers.AdminDashboardHandler
}

func NewAppContainer() *AppContainer {
	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	teamsRepo := repositories.NewTeamsRepository()
	teamMemberRepo := repositories.NewTeamMemberRepository()

	// Initialize services
	authService := services.NewAuthService(config.DB, userRepo)
	userService := services.NewUserService(config.DB, userRepo)
	teamsService := services.NewTeamsService(config.DB, teamsRepo, teamMemberRepo)

	return &AppContainer{
		// Middlewares
		JWTAuthMiddleware:   middlewares.JWTAuthMiddleware(),
		AdminAuthMiddleware: middlewares.AdminAuthMiddleware(),

		// Services
		AuthService: authService,
		UserService: userService,
		TeamsService: teamsService,

		// Handlers
		AuthHandler:        handlers.NewAuthHandler(authService),
		DashboardHandler:   handlers.NewDashboardHandler(),
		UserProfileHandler: handlers.NewUserProfileHandler(userService),
		TeamsHandler:       handlers.NewTeamsHandler(teamsService),
		// Admin Handlers
		AdminAuthHandler:      handlers.NewAdminAuthHandler(authService),
		AdminDashboardHandler: handlers.NewAdminDashboardHandler(),
	}
}
