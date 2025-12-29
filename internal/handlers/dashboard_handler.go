package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) DashboardPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/dashboard.html", gin.H{
		"title": "Dashboard",
	})
}
