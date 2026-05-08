package discord

import (
	"context"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/dto"
)

type DiscordOAuth interface {
	GetAccessToken(ctx context.Context, clientID string, clientSecret string, code string, redirectURI string) (string, error)
	GetGuildMemberData(ctx context.Context, accessToken string, guildID string) (*dto.DiscordGuildMemberDTO, error)
	GetConnections(ctx context.Context, accessToken string) (string, error)
}
