package handlers

import (
	"net/http"
	"trieu_mock_project_go/internal/services"
	"trieu_mock_project_go/internal/types"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type AdminAuthHandler struct {
	activityLogService *services.ActivityLogService
	authService        *services.AuthService
}

func NewAdminAuthHandler(authService *services.AuthService, activityLogService *services.ActivityLogService) *AdminAuthHandler {
	return &AdminAuthHandler{authService: authService, activityLogService: activityLogService}
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

	user, err := h.authService.Login(c.Request.Context(), email, password, true)
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
	session.Set("email", user.Email)
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

	userIdValue := session.Get("user_id")
	emailValue := session.Get("email")

	userId, getUserIdOk := userIdValue.(uint)
	email, getEmailOk := emailValue.(string)

	if !getUserIdOk || !getEmailOk {
		session.Clear()
		_ = session.Save()

		c.Redirect(http.StatusSeeOther, "/admin/login")
		return
	}

	// Log activity
	if err := h.activityLogService.LogActivity(c.Request.Context(), types.AdminSignOut, userId, email); err != nil {
		c.HTML(http.StatusInternalServerError, "pages/admin_login.html", gin.H{
			"title":     "Admin Login",
			"error":     "Failed to log activity",
			"csrfToken": csrf.GetToken(c),
		})
		return
	}

	// Clear session
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
