package handlers

import (
	"net/http"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminAuthHandler struct {
	authService *services.AuthService
}

func NewAdminAuthHandler(authService *services.AuthService) *AdminAuthHandler {
	return &AdminAuthHandler{authService: authService}
}

func (h *AdminAuthHandler) AdminShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_login.html", gin.H{
		"title": "Admin Login",
	})
}

func (h *AdminAuthHandler) AdminLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := h.authService.Login(c.Request.Context(), email, password)
	if err != nil || user.Role != "admin" {
		c.HTML(http.StatusUnauthorized, "pages/login.html", gin.H{
			"title": "Admin Login",
			"error": "Invalid email or password, or not an admin",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("role", user.Role)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/admin")
}

func (h *AdminAuthHandler) AdminLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}
