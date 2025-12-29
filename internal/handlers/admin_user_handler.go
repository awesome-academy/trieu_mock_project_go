package handlers

import (
	"net/http"
	"strconv"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	userService     *services.UserService
	teamService     *services.TeamsService
	positionService *services.PositionService
	skillService    *services.SkillService
}

func NewAdminUserHandler(
	userService *services.UserService,
	teamService *services.TeamsService,
	positionService *services.PositionService,
	skillService *services.SkillService) *AdminUserHandler {
	return &AdminUserHandler{
		userService:     userService,
		teamService:     teamService,
		positionService: positionService,
		skillService:    skillService,
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

func (h *AdminUserHandler) AdminUserEditPage(c *gin.Context) {
	userIdParam := c.Param("userId")

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.HTML(http.StatusBadRequest, "pages/admin_user_edit.html", gin.H{
			"title": "Edit User",
			"error": "Invalid user ID",
		})
		return
	}

	userProfile, err := h.userService.GetUserProfile(c.Request.Context(), uint(userId))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "pages/admin_user_edit.html", gin.H{
			"title": "Edit User",
			"error": "Failed to load user details",
		})
		return
	}

	allTeam := h.teamService.GetAllTeamsSummary(c.Request.Context())
	positions := h.positionService.GetAllPositionsSummary(c.Request.Context())
	skills := h.skillService.GetAllSkillsSummary(c.Request.Context())

	c.HTML(http.StatusOK, "pages/admin_user_edit.html", gin.H{
		"title":     "Edit User",
		"user":      userProfile,
		"teams":     allTeam,
		"positions": positions,
		"skills":    skills,
	})
}

func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	userIdParam := c.Param("userId")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dtos.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.UpdateUser(c.Request.Context(), uint(userId), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
