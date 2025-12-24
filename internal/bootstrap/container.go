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
	RequireLoginMiddleware gin.HandlerFunc
	RequireAdminMiddleware gin.HandlerFunc

	// Services
	UserService *services.UserService

	// Handlers
	ShowLoginHandler     gin.HandlerFunc
	DoLoginHandler       gin.HandlerFunc
	LogoutHandler        gin.HandlerFunc
	DashboardPageHandler gin.HandlerFunc
	// Admin Handlers
	AdminDashboardPageHandler gin.HandlerFunc
}

func NewAppContainer() *AppContainer {
	// Initialize repositories
	userRepo := repositories.NewUserRepository()

	// Initialize services
	userService := services.NewUserService(config.DB, userRepo)

	return &AppContainer{
		// Middlewares
		RequireLoginMiddleware: middlewares.RequireLogin(),
		RequireAdminMiddleware: middlewares.RequireAdmin(),

		// Services
		UserService: userService,

		// Handlers
		ShowLoginHandler:     handlers.ShowLoginHandler,
		DoLoginHandler:       handlers.NewDoLoginHandler(userService),
		LogoutHandler:        handlers.LogoutHandler,
		DashboardPageHandler: handlers.DashboardPageHandler,
		// Admin Handlers
		AdminDashboardPageHandler: handlers.AdminDashboardPageHandler,
	}
}
