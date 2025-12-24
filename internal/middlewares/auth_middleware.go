package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_id") == nil {
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Redirect", "/login")
			} else {
				c.Redirect(302, "/login")
			}
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")

		if role != "admin" {
			// Redirect to home page if not admin
			c.Redirect(302, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}
