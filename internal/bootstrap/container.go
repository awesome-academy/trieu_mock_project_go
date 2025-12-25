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

	// Handlers
	AuthHandler        *handlers.AuthHandler
	DashboardHandler   *handlers.DashboardHandler
	UserProfileHandler *handlers.UserProfileHandler
	// Admin Handlers
	AdminAuthHandler      *handlers.AdminAuthHandler
	AdminDashboardHandler *handlers.AdminDashboardHandler
}

func NewAppContainer() *AppContainer {
	// Initialize repositories
	userRepo := repositories.NewUserRepository()

	// Initialize services
	authService := services.NewAuthService(config.DB, userRepo)
	userService := services.NewUserService(config.DB, userRepo)

	return &AppContainer{
		// Middlewares
		JWTAuthMiddleware:   middlewares.JWTAuthMiddleware(),
		AdminAuthMiddleware: middlewares.AdminAuthMiddleware(),

		// Services
		AuthService: authService,
		UserService: userService,

		// Handlers
		AuthHandler:        handlers.NewAuthHandler(authService),
		DashboardHandler:   handlers.NewDashboardHandler(),
		UserProfileHandler: handlers.NewUserProfileHandler(userService),
		// Admin Handlers
		AdminAuthHandler:      handlers.NewAdminAuthHandler(authService),
		AdminDashboardHandler: handlers.NewAdminDashboardHandler(),
	}
}
