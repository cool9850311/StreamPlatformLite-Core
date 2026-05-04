package middleware

import (
	"context"
	"strings"

	coreClaims "github.com/cool9850311/StreamPlatformLite-Core/pkg/claims"
	"github.com/cool9850311/StreamPlatformLite-Core/pkg/role"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/interface/logger"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// OptionalJWTAuthMiddleware is an optional JWT auth middleware.
// If a valid token is present, parse it and store in context.
// If no token or invalid token, set to Anonymous role.
func OptionalJWTAuthMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Try reading from Cookie
		cookie, err := c.Cookie("token")
		if err == nil && cookie != "" {
			tokenString = cookie
		} else {
			// 2. Try reading from Authorization Header
			tokenString = c.GetHeader("Authorization")
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		// If no token, set to Anonymous and continue
		if tokenString == "" {
			claims := &coreClaims.Claims{
				Role: role.Anonymous,
			}
			ctx := context.WithValue(c.Request.Context(), "claims", claims)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			return
		}

		// 3. Try parsing JWT
		claims := &coreClaims.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWT.SecretKey), nil
		})

		// If token is invalid, set to Anonymous
		if err != nil || !token.Valid {
			logger.Warn(c.Request.Context(), "Invalid JWT token, treating as anonymous: "+err.Error())
			claims = &coreClaims.Claims{
				Role: role.Anonymous,
			}
			ctx := context.WithValue(c.Request.Context(), "claims", claims)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			return
		}

		// Token is valid, store in context
		ctx := context.WithValue(c.Request.Context(), "claims", claims)
		c.Request = c.Request.WithContext(ctx)

		// Also store user_id in Gin context for rate limiting middleware
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
