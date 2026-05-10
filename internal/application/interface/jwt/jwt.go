package jwt

import (
	"context"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/dto"
	"github.com/cool9850311/StreamPlatformLite-Core/pkg/role"
)

type JWTGenerator interface {
	GenerateDiscordToken(ctx context.Context, discordId string, guildMemberData *dto.DiscordGuildMemberDTO, userRole role.Role, secretKey string, ytChannelID string, ytName string) (string, error)
	GenerateOriginToken(ctx context.Context, userID string, username string, userRole role.Role, secretKey string) (string, error)
}
