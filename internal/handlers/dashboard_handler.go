package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/dashboard.html", gin.H{
		"title":      "Home",
	})
}
