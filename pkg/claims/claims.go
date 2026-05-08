package claims

import (
	"github.com/cool9850311/StreamPlatformLite-Core/pkg/role"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID           string    `json:"user_id"`
	Avatar           string    `json:"avatar"`
	UserName         string    `json:"username"`
	Role             role.Role `json:"role"`
	IdentityProvider string    `json:"identity_provider"`
	YtChannelID      string    `json:"yt_channel_id"`
	jwt.RegisteredClaims
}
