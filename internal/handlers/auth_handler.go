package handlers

import (
	"net/http"
	"time"
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
	user, err := h.authService.Login(c.Request.Context(), req.User.Email, req.User.Password, false)
	if err != nil {
		appErrors.RespondError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Email)
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	if err := h.authService.StoreToken(c.Request.Context(), user.ID, token, 24*time.Hour); err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to store session")
		return
	}

	resp := dtos.LoginResponse{}
	resp.User.ID = user.ID
	resp.User.Name = user.Name
	resp.User.Email = user.Email
	resp.User.AccessToken = token

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) GenerateWSTicket(c *gin.Context) {
	userID := c.GetUint("user_id")
	email := c.GetString("email")

	ticket, err := h.authService.CreateWSTicket(c.Request.Context(), userID, email)
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to generate WebSocket ticket")
		return
	}

	c.JSON(http.StatusOK, dtos.WSTicketResponse{Ticket: ticket})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetUint("user_id")
	token := c.GetString("token")
	email := c.GetString("user_email")

	if userID != 0 && token != "" {
		if err := h.authService.Logout(c.Request.Context(), userID, token, email); err != nil {
			appErrors.RespondCustomError(c, err, "Failed to logout")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
