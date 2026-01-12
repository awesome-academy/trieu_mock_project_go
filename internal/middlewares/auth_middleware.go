package middlewares

import (
	"context"
	"net/http"
	"strings"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"
	"trieu_mock_project_go/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// extractAndValidateToken extracts and validates JWT token from Authorization header
// Returns (userID, email, tokenString, error) with error constants from errors package
func extractAndValidateToken(c *gin.Context) (uint, string, string, error) {
	authHeader := c.GetHeader("Authorization")
	tokenString := ""

	if authHeader == "" {
		return 0, "", "", appErrors.ErrMissingAuthHeader
	}

	// Extract token from "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		tokenString = parts[1]
	} else {
		return 0, "", "", appErrors.ErrInvalidAuthHeader
	}

	if tokenString == "" {
		return 0, "", "", appErrors.ErrMissingAuthHeader
	}

	claims, err := utils.ParseJWTToken(tokenString)
	if err != nil {
		return 0, "", "", appErrors.ErrInvalidToken
	}

	return claims.UserID, claims.Email, tokenString, nil
}

// JWTAuthMiddleware checks JWT token from Authorization header and verifies it in Redis
func JWTAuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, email, tokenString, err := extractAndValidateToken(c)
		if err != nil {
			switch err {
			case appErrors.ErrMissingAuthHeader, appErrors.ErrInvalidAuthHeader, appErrors.ErrInvalidToken:
				appErrors.RespondError(c, http.StatusUnauthorized, err.Error())
			default:
				appErrors.RespondError(c, http.StatusUnauthorized, "authentication failed")
			}
			c.Abort()
			return
		}

		// Check if token exists in Redis
		isValid, err := authService.IsTokenStoreValid(c.Request.Context(), userID, tokenString)
		if err != nil || !isValid {
			appErrors.RespondError(c, http.StatusUnauthorized, "token is invalid or expired")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("token", tokenString)
		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_id") == nil {
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Redirect", "/admin/login")
				c.Abort()
			} else {
				c.Redirect(302, "/admin/login")
				c.Abort()
			}
			return
		}
		role := session.Get("role")

		if role != "admin" {
			c.Redirect(http.StatusForbidden, "/forbidden")
			c.Abort()
			return
		}

		// Set user info to context for logging activities
		ctx := c.Request.Context()
		if userID := session.Get("user_id"); userID != nil {
			ctx = context.WithValue(ctx, "user_id", userID)
		}
		if email := session.Get("email"); email != nil {
			ctx = context.WithValue(ctx, "email", email)
		}
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}

}

func JWTAuthWSMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ticket := c.Query("ticket")
		if ticket == "" {
			appErrors.RespondError(c, http.StatusUnauthorized, "missing websocket ticket")
			c.Abort()
			return
		}

		userID, email, err := authService.ConsumeWSTicket(c.Request.Context(), ticket)
		if err != nil {
			appErrors.RespondCustomError(c, err, "invalid or expired websocket ticket")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Next()
	}
}
