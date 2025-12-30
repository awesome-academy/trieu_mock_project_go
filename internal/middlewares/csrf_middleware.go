package middlewares

import (
	"net/http"
	"trieu_mock_project_go/internal/config"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func CSRFMiddleware() gin.HandlerFunc {
	cfg := config.LoadConfig()
	return csrf.Middleware(csrf.Options{
		Secret: cfg.SessionConfig.Secret,
		ErrorFunc: func(c *gin.Context) {
			c.String(http.StatusBadRequest, "CSRF token mismatch")
			c.Abort()
		},
	})
}
