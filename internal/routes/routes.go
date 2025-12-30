package routes

import (
	"trieu_mock_project_go/internal/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, appContainer *bootstrap.AppContainer) {
	// User endpoint
	router.GET("/login", appContainer.AuthHandler.ShowLoginPage)
	router.POST("/login", appContainer.AuthHandler.UserLogin)
	router.GET("/", appContainer.DashboardHandler.DashboardPageHandler)
	router.GET("/profile", appContainer.UserProfileHandler.UserProfilePageHandler)
	router.GET("/teams", appContainer.TeamsHandler.TeamsPageHandler)

	// Normal user routes (JWT)
	apiGroup := router.Group("/api")
	apiGroup.Use(appContainer.JWTAuthMiddleware)
	{
		apiGroup.GET("/profile", appContainer.UserProfileHandler.GetUserProfile)
		apiGroup.GET("/teams", appContainer.TeamsHandler.ListTeams)
	}

	// Admin login flow
	router.GET("/admin/login", appContainer.AdminAuthHandler.AdminShowLogin)
	router.POST("/admin/login", appContainer.AdminAuthHandler.AdminLogin)
	router.GET("/admin/logout", appContainer.AdminAuthHandler.AdminLogout)

	// Admin routes (Session)
	adminGroup := router.Group("/admin")
	adminGroup.Use(appContainer.AdminAuthMiddleware)
	{
		adminGroup.GET("/", appContainer.AdminDashboardHandler.AdminDashboardPage)
	}
}
