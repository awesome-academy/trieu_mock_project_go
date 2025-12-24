package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.page.tmpl", gin.H{
		"title": "Login Page",
	})
}

func DoLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	session := sessions.Default(c)
	if email == "admin@sun-asterisk.com" && password == "password" {
		session.Set("user_id", 1)
		session.Set("role", "admin")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	if email == "user@sun-asterisk.com" && password == "password" {
		session.Set("user_id", 2)
		session.Set("role", "user")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	c.HTML(http.StatusUnauthorized, "login.page.tmpl", gin.H{
		"title": "Login Page",
		"error": "Invalid email or password",
	})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}
