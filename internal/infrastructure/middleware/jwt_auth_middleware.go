package middleware

import (
	"context"
	"net/http"

	coreClaims "github.com/cool9850311/StreamPlatformLite-Core/pkg/claims"
	"github.com/cool9850311/StreamPlatformLite-Core/pkg/csrf"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/interface/logger"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware handles JWT authentication and authorization
func JWTAuthMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
			c.Abort()
			return
		}

		claims := &coreClaims.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWT.SecretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		// Store claims in context for later use
		ctx := context.WithValue(c.Request.Context(), "claims", claims)
		c.Request = c.Request.WithContext(ctx)

		// Also store user_id in Gin context for rate limiting middleware
		c.Set("user_id", claims.UserID)

		// CSRF: enforce on all state-changing methods
		safeMethod := c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS"
		if !safeMethod {
			csrfToken := c.GetHeader("X-XSRF-TOKEN")
			if csrfToken == "" || !csrf.ValidateCsrfToken(csrfToken, config.AppConfig.JWT.SecretKey, claims.UserID) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "invalid CSRF token"})
				return
			}
		}

		c.Next()
	}
}
