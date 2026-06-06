package mock_data

import (
	"context"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/dto"
	"github.com/cool9850311/StreamPlatformLite-Core/pkg/role"
	"github.com/stretchr/testify/mock"
)

type MockJWTGenerator struct {
	mock.Mock
}

func (m *MockJWTGenerator) GenerateDiscordToken(ctx context.Context, discordId string, guildMemberData *dto.DiscordGuildMemberDTO, userRole role.Role, secretKey string, ytChannelID string, ytName string, roleIDs []string) (string, error) {
	args := m.Called(ctx, discordId, guildMemberData, userRole, secretKey, ytChannelID, ytName, roleIDs)
	return args.String(0), args.Error(1)
}

func (m *MockJWTGenerator) GenerateOriginToken(ctx context.Context, userID string, username string, userRole role.Role, secretKey string) (string, error) {
	args := m.Called(ctx, userID, username, userRole, secretKey)
	return args.String(0), args.Error(1)
}
