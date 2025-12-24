package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/dashboard.html", gin.H{
		"title": "Dashboard",
	})
}
