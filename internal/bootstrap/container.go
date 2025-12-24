package bootstrap

import (
	"trieu_mock_project_go/internal/handlers"
	"trieu_mock_project_go/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type AppContainer struct {
	// Middlewares
	RequireLoginMiddleware gin.HandlerFunc
	RequireAdminMiddleware gin.HandlerFunc

	// Handlers
	ShowLoginHandler     gin.HandlerFunc
	DoLoginHandler       gin.HandlerFunc
	LogoutHandler        gin.HandlerFunc
	ShowDashboardHandler gin.HandlerFunc
	// Admin Handlers
	ShowAdminDashboardHandler gin.HandlerFunc
}

func NewAppContainer() *AppContainer {
	return &AppContainer{
		// Middlewares
		RequireLoginMiddleware: middlewares.RequireLogin(),
		RequireAdminMiddleware: middlewares.RequireAdmin(),

		// Handlers
		ShowLoginHandler:     handlers.ShowLogin,
		DoLoginHandler:       handlers.DoLogin,
		LogoutHandler:        handlers.Logout,
		ShowDashboardHandler: handlers.ShowDashboard,
		// Admin Handlers
		ShowAdminDashboardHandler: handlers.ShowAdminDashboard,
	}
}
