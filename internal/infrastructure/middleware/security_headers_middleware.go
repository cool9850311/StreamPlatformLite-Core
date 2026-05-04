package middleware

import (
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/config"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func SecurityHeaders() gin.HandlerFunc {
	// HSTS: only enable in production HTTPS mode (STSSeconds > 0 adds HSTS header)
	stsSeconds := int64(0)
	if config.AppConfig.Server.HTTPS {
		stsSeconds = 31536000 // 1 year (31536000 seconds)
	}

	secureConfig := secure.Config{
		// X-Frame-Options: DENY (prevents Clickjacking)
		FrameDeny: true,

		// X-Content-Type-Options: nosniff (prevents MIME sniffing)
		ContentTypeNosniff: true,

		// X-XSS-Protection: 1; mode=block
		BrowserXssFilter: true,

		// Referrer-Policy: strict-origin-when-cross-origin
		ReferrerPolicy: "strict-origin-when-cross-origin",

		// HSTS: only set when HTTPS is enabled (STSSeconds > 0)
		STSSeconds:           stsSeconds,
		STSIncludeSubdomains: true,
		STSPreload:           false,

		IsDevelopment: false,

		// Zero-trust API CSP: prevents API responses from being treated as HTML
		// and prevents embedding in iframes
		ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
	}

	secureMiddleware := secure.New(secureConfig)

	// Return a wrapper that adds Permissions-Policy header manually
	return func(c *gin.Context) {
		// Add Permissions-Policy header (newer standard replacing Feature-Policy)
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Call the secure middleware
		secureMiddleware(c)
	}
}
