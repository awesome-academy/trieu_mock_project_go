package handlers

import (
	"net/http"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type AdminAuthHandler struct {
	authService *services.AuthService
}

func NewAdminAuthHandler(authService *services.AuthService) *AdminAuthHandler {
	return &AdminAuthHandler{authService: authService}
}

func (h *AdminAuthHandler) AdminShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_login.html", gin.H{
		"title":     "Admin Login",
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminAuthHandler) AdminLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := h.authService.Login(c.Request.Context(), email, password)
	if err != nil || user.Role != "admin" {
		c.HTML(http.StatusUnauthorized, "pages/admin_login.html", gin.H{
			"title":     "Admin Login",
			"error":     "Invalid email or password, or not an admin",
			"csrfToken": csrf.GetToken(c),
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("role", user.Role)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "pages/admin_login.html", gin.H{
			"title":     "Admin Login",
			"error":     "Failed to save session",
			"csrfToken": csrf.GetToken(c),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin")
}

func (h *AdminAuthHandler) AdminLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "pages/admin_login.html", gin.H{
			"title":     "Admin Login",
			"error":     "Failed to save session",
			"csrfToken": csrf.GetToken(c),
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/login")
}
