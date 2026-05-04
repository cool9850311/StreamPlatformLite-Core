package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

// RateLimitByIP - IP dimension rate limiting middleware
func RateLimitByIP(lim *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.AppConfig.RateLimit.Enabled {
			c.Next()
			return
		}

		// Get client IP (prefer X-Forwarded-For)
		ip := c.GetHeader("X-Forwarded-For")
		if ip == "" {
			ip = c.ClientIP()
		}

		limiterCtx, err := lim.Get(context.Background(), ip)
		if err != nil {
			// Limiter error, fail-open (allow request through)
			c.Next()
			return
		}

		// Set response headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

		if limiterCtx.Reached {
			retryAfterSeconds := limiterCtx.Reset - time.Now().Unix()
			if retryAfterSeconds < 0 {
				retryAfterSeconds = 0
			}
			c.Header("Retry-After", strconv.FormatInt(retryAfterSeconds, 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message":     "Rate limit exceeded. Please try again later.",
				"retry_after": retryAfterSeconds,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUserID - user ID dimension rate limiting middleware
func RateLimitByUserID(lim *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.AppConfig.RateLimit.Enabled {
			c.Next()
			return
		}

		// Get user ID from JWT claims
		userID, exists := c.Get("user_id")
		if !exists {
			// Unauthenticated user, skip rate limiting
			c.Next()
			return
		}

		key := fmt.Sprintf("user:%v", userID)
		limiterCtx, err := lim.Get(context.Background(), key)
		if err != nil {
			// fail-open
			c.Next()
			return
		}

		// Set response headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

		if limiterCtx.Reached {
			retryAfterSeconds := limiterCtx.Reset - time.Now().Unix()
			if retryAfterSeconds < 0 {
				retryAfterSeconds = 0
			}
			c.Header("Retry-After", strconv.FormatInt(retryAfterSeconds, 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message":     "Rate limit exceeded. Please try again later.",
				"retry_after": retryAfterSeconds,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
