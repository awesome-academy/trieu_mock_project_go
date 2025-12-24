package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "home.page.tmpl", gin.H{
		"title":      "Home",
	})
}
