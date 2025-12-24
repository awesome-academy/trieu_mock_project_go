package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowAdminDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_dashboard.html", gin.H{
		"title":      "Admin Dashboard",
	})
}
