package handlers

import (
	"net/http"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"
	"trieu_mock_project_go/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/login.html", gin.H{
		"title": "User Login",
	})
}

func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req dtos.LoginRequest

	// Validate request body
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&req)) {
		return
	}

	// Login user
	user, err := h.authService.Login(c.Request.Context(), req.User.Email, req.User.Password)
	if err != nil {
		appErrors.RespondError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Email)
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	resp := dtos.LoginResponse{}
	resp.User.ID = user.ID
	resp.User.Name = user.Name
	resp.User.Email = user.Email
	resp.User.AccessToken = token

	c.JSON(http.StatusOK, resp)
}
