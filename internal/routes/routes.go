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
	router.GET("/admin/login", appContainer.AdminAuthHandler.AdminShowLogin)
	router.POST("/admin/login", appContainer.AdminAuthHandler.AdminLogin)
	router.GET("/admin/logout", appContainer.AdminAuthHandler.AdminLogout)

	// Admin routes (Session)
	adminGroup := router.Group("/admin")
	adminGroup.Use(appContainer.AdminAuthMiddleware)
	{
		adminGroup.GET("/", appContainer.AdminDashboardHandler.AdminDashboardPage)
		adminGroup.GET("/users", appContainer.AdminUserHandler.AdminUsersPage)
		adminGroup.GET("/users/partial/search", appContainer.AdminUserHandler.AdminUsersSearchPartial)
		adminGroup.GET("/users/create", appContainer.AdminUserHandler.AdminUserCreatePage)
		adminGroup.POST("/users", appContainer.AdminUserHandler.CreateUser)
		adminGroup.GET("/users/:userId", appContainer.AdminUserHandler.AdminUserDetailPage)
		adminGroup.GET("/users/:userId/edit", appContainer.AdminUserHandler.AdminUserEditPage)
		adminGroup.PUT("/users/:userId", appContainer.AdminUserHandler.UpdateUser)
		adminGroup.DELETE("/users/:userId", appContainer.AdminUserHandler.DeleteUser)
	}
}
