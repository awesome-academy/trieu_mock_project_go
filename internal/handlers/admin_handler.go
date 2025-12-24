package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowAdminDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_dashboard.page.tmpl", gin.H{
		"title":      "Admin Dashboard",
	})
}
