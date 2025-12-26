package handlers

import (
	"net/http"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminDashboardHandler struct {
	userService *services.UserService
}

func NewAdminDashboardHandler(userService *services.UserService) *AdminDashboardHandler {
	return &AdminDashboardHandler{
		userService: userService,
	}
}

func (h *AdminDashboardHandler) AdminDashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_dashboard.html", gin.H{
		"title": "Admin Dashboard",
	})
}