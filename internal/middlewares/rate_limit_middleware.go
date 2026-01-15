package middlewares

import (
	"fmt"
	"net/http"
	"time"
	"trieu_mock_project_go/internal/config"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(redisService *services.RedisService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.LoadConfig()
		ip := c.ClientIP()
		now := time.Now().Unix()
		key := fmt.Sprintf("ratelimit:%s:%d", ip, now)

		count, err := redisService.Incr(c.Request.Context(), key)
		if err != nil {
			// In case of Redis error, we allow the request to proceed to avoid breaking the app
			c.Next()
			return
		}

		if count == 1 {
			if _, err := redisService.Expire(c.Request.Context(), key, 1*time.Second); err != nil {
				fmt.Printf("failed to set expiration for rate limit key %s: %v\n", key, err)
			}
		}

		if count > int64(cfg.RequestRateLimit) {
			appErrors.RespondError(c, http.StatusTooManyRequests, appErrors.ErrTooManyRequests.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
