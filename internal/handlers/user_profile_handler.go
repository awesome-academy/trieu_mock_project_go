package handlers

import (
	"net/http"
	"strconv"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	userService *services.UserService
}

func NewUserProfileHandler(userService *services.UserService) *UserProfileHandler {
	return &UserProfileHandler{
		userService: userService,
	}
}

func (h *UserProfileHandler) UserMyProfilePageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/user_profile.html", gin.H{
		"title": "User Profile",
	})
}

func (h *UserProfileHandler) UserUserProfilePageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/user_profile.html", gin.H{
		"title": "User Profile",
	})
}

func (h *UserProfileHandler) GetMyProfile(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		appErrors.RespondError(c, http.StatusUnauthorized, "Unauthorized access")
		return
	}
	userProfile, err := h.userService.GetUserProfile(c.Request.Context(), userId)
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

func (h *UserProfileHandler) GetUserProfile(c *gin.Context) {
	userIdParam := c.Param("userId")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userProfile, err := h.userService.GetUserProfile(c.Request.Context(), uint(userId))
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	c.JSON(http.StatusOK, userProfile)
}
