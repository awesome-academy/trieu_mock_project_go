package middlewares

import (
	"fmt"
	"net/http"
	"time"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

// Rate limiting: max 5 requests/s/IP
const RateLimitRequests = 5

func RateLimitMiddleware(redisService *services.RedisService) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			redisService.Expire(c.Request.Context(), key, 2*time.Second)
		}

		if count > RateLimitRequests {
			appErrors.RespondError(c, http.StatusTooManyRequests, appErrors.ErrTooManyRequests.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
