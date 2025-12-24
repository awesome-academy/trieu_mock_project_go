package routes

import (
	"trieu_mock_project_go/internal/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, appContainer *bootstrap.AppContainer) {
	// public endpoint
	router.GET("/login", appContainer.ShowLoginHandler)
	router.POST("/login", appContainer.DoLoginHandler)
	router.GET("/logout", appContainer.LogoutHandler)

	router.Use(appContainer.RequireLoginMiddleware)
	router.GET("/", appContainer.ShowDashboardHandler)

	adminGroup := router.Group("/admin")
	adminGroup.Use(appContainer.RequireAdminMiddleware)
	{
		adminGroup.GET("/", appContainer.ShowAdminDashboardHandler)
	}
}
