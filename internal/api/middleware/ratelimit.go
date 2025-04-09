package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/0xBoji/web3-edu-core/internal/database/redis"
	"github.com/0xBoji/web3-edu-core/internal/utils"
	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware is a middleware for rate limiting
func RateLimitMiddleware(requests int64, duration time.Duration) gin.HandlerFunc {
	cache := redis.NewCache()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		ctx := context.Background()

		// Increment the rate limit counter
		count, err := cache.IncrementRateLimit(ctx, ip, duration)
		if err != nil {
			utils.ServerErrorResponse(c)
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", "100")
		c.Header("X-RateLimit-Reset", "60")

		// Check if the rate limit has been exceeded
		if count > requests {
			c.Header("Retry-After", "60")
			utils.ErrorResponse(c, http.StatusTooManyRequests, "rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
