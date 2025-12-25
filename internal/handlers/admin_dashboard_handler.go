package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminDashboardHandler struct{}

func NewAdminDashboardHandler() *AdminDashboardHandler {
	return &AdminDashboardHandler{}
}

func (h *AdminDashboardHandler) AdminDashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_dashboard.html", gin.H{
		"title": "Admin Dashboard",
	})
}
