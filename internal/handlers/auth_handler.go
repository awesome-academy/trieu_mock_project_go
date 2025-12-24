package handlers

import (
	"net/http"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ShowLoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/login.html", gin.H{
		"title": "Login Page",
	})
}

func NewDoLoginHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		user, err := userService.Login(email, password)
		if err != nil {
			c.HTML(http.StatusUnauthorized, "pages/login.html", gin.H{
				"title": "Login Page",
				"error": "Invalid email or password",
			})
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", user.ID)
		session.Set("role", user.Role)
		session.Save()

		if user.Role == "admin" {
			c.Redirect(http.StatusSeeOther, "/admin")
			return
		}

		c.Redirect(http.StatusSeeOther, "/")
	}
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}
