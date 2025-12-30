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
	router.GET("/profile", appContainer.UserProfileHandler.UserMyProfilePageHandler)
	router.GET("/profile/:userId", appContainer.UserProfileHandler.UserUserProfilePageHandler)
	router.GET("/teams", appContainer.TeamsHandler.TeamsPageHandler)
	router.GET("/teams/:id", appContainer.TeamsHandler.TeamDetailsPageHandler)

	// Normal user routes (JWT)
	apiGroup := router.Group("/api")
	apiGroup.Use(appContainer.JWTAuthMiddleware)
	{
		apiGroup.GET("/profile", appContainer.UserProfileHandler.GetMyProfile)
		apiGroup.GET("/profile/:userId", appContainer.UserProfileHandler.GetUserProfile)
		apiGroup.GET("/teams", appContainer.TeamsHandler.ListTeams)
		apiGroup.GET("/teams/:id", appContainer.TeamsHandler.GetTeamDetails)
		apiGroup.GET("/teams/:id/members", appContainer.TeamsHandler.GetTeamMembers)
	}

	// Admin login flow
	router.GET("/admin/login", appContainer.CSRFMiddleware, appContainer.AdminAuthHandler.AdminShowLogin)
	router.POST("/admin/login", appContainer.CSRFMiddleware, appContainer.AdminAuthHandler.AdminLogin)
	router.GET("/admin/logout", appContainer.AdminAuthHandler.AdminLogout)

	// Admin routes (Session)
	adminGroup := router.Group("/admin")
	adminGroup.Use(appContainer.AdminAuthMiddleware)
	{
		// Admin dashboard
		adminGroup.GET("/", appContainer.AdminDashboardHandler.AdminDashboardPage)
		// Admin user management
		adminGroup.GET("/users", appContainer.AdminUserHandler.AdminUsersPage)
		adminGroup.GET("/users/partial/search", appContainer.AdminUserHandler.AdminUsersSearchPartial)
		adminGroup.GET("/users/create", appContainer.CSRFMiddleware, appContainer.AdminUserHandler.AdminUserCreatePage)
		adminGroup.POST("/users", appContainer.CSRFMiddleware, appContainer.AdminUserHandler.CreateUser)
		adminGroup.GET("/users/:userId", appContainer.CSRFMiddleware, appContainer.AdminUserHandler.AdminUserDetailPage)
		adminGroup.GET("/users/:userId/edit", appContainer.CSRFMiddleware, appContainer.AdminUserHandler.AdminUserEditPage)
		adminGroup.PUT("/users/:userId", appContainer.CSRFMiddleware, appContainer.AdminUserHandler.UpdateUser)
		adminGroup.DELETE("/users/:userId", appContainer.CSRFMiddleware, appContainer.AdminUserHandler.DeleteUser)
		// Admin position management
		adminGroup.GET("/positions", appContainer.AdminPositionHandler.ListPositionPage)
		adminGroup.GET("/positions/partial/search", appContainer.AdminPositionHandler.PositionSearchPartial)
		adminGroup.GET("/positions/create", appContainer.AdminPositionHandler.CreatePositionPage)
		adminGroup.POST("/positions", appContainer.AdminPositionHandler.CreatePosition)
		adminGroup.GET("/positions/:positionId/edit", appContainer.AdminPositionHandler.EditPositionPage)
		adminGroup.PUT("/positions/:positionId", appContainer.AdminPositionHandler.UpdatePosition)
		adminGroup.DELETE("/positions/:positionId", appContainer.AdminPositionHandler.DeletePosition)
	}
}
