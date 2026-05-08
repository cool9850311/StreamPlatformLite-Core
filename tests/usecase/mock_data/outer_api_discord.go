package mock_data

import (
	"context"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/dto"
	"github.com/stretchr/testify/mock"
)

type MockDiscordOAuth struct {
	mock.Mock
}

// Mock implementation of GetAccessToken
func (m *MockDiscordOAuth) GetAccessToken(ctx context.Context, clientID string, clientSecret string, code string, redirectURI string) (string, error) {
	args := m.Called(ctx, clientID, clientSecret, code, redirectURI)
	return args.String(0), args.Error(1)
}

// Mock implementation of GetGuildMemberData
func (m *MockDiscordOAuth) GetGuildMemberData(ctx context.Context, accessToken string, guildID string) (*dto.DiscordGuildMemberDTO, error) {
	args := m.Called(ctx, accessToken, guildID)
	return args.Get(0).(*dto.DiscordGuildMemberDTO), args.Error(1)
}

// Mock implementation of GetConnections
func (m *MockDiscordOAuth) GetConnections(ctx context.Context, accessToken string) (string, error) {
	args := m.Called(ctx, accessToken)
	return args.String(0), args.Error(1)
}
