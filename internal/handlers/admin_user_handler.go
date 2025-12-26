package handlers

import (
	"net/http"
	"strconv"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	userService *services.UserService
	teamService *services.TeamsService
}

func NewAdminUserHandler(userService *services.UserService, teamService *services.TeamsService) *AdminUserHandler {
	return &AdminUserHandler{
		userService: userService,
		teamService: teamService,
	}
}

func (h *AdminUserHandler) AdminUsersPage(c *gin.Context) {
	allTeam := h.teamService.GetAllTeamsSummary(c.Request.Context())
	c.HTML(http.StatusOK, "pages/admin_users.html", gin.H{
		"title": "Admin Users Management",
		"teams": allTeam,
	})
}

func (h *AdminUserHandler) AdminUsersSearchPartial(c *gin.Context) {
	var query dtos.UserSearchRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.HTML(http.StatusBadRequest, "partials/admin_users_search.html", gin.H{
			"error": "Invalid query parameters",
		})
		return
	}

	resp, err := h.userService.SearchUsers(c.Request.Context(), query.TeamId, query.Limit, query.Offset)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "partials/admin_users_search.html", gin.H{
			"error": "Failed to load users",
		})
		return
	}

	c.HTML(http.StatusOK, "partials/admin_users_search.html", gin.H{
		"users": resp.Users,
		"page":  resp.Page,
	})
}

func (h *AdminUserHandler) AdminUserDetailPage(c *gin.Context) {
	userIdParam := c.Param("userId")

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.HTML(http.StatusBadRequest, "pages/admin_user_detail.html", gin.H{
			"title": "User Detail",
			"error": "Invalid user ID",
		})
		return
	}

	userProfile, err := h.userService.GetUserProfile(c.Request.Context(), uint(userId))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "pages/admin_user_detail.html", gin.H{
			"title": "User Detail",
			"error": "Failed to load user details",
		})
		return
	}

	c.HTML(http.StatusOK, "pages/admin_user_detail.html", gin.H{
		"title": "User Detail",
		"user":  userProfile,
	})
}
