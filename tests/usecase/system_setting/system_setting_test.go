package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/usecase"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/entity/system"
	"github.com/cool9850311/StreamPlatformLite-Core/pkg/role"
	"github.com/cool9850311/StreamPlatformLite-Core/tests/usecase/mock_data"
	"github.com/stretchr/testify/assert"
)

func setup() (*mock_data.MockSystemSettingRepository, *mock_data.MockLogger, usecase.SystemSettingUseCase) {
	mockRepo := new(mock_data.MockSystemSettingRepository)
	mockLogger := new(mock_data.MockLogger)

	useCase := usecase.NewSystemSettingUseCase(mockRepo, mockLogger)
	return mockRepo, mockLogger, *useCase
}

func TestSystemSettingUseCase_GetSetting_AdminUser(t *testing.T) {
	mockRepo, _, useCase := setup()
	ctx := context.Background()

	testSetting := &system.Setting{
		EditorRoleId:        "editor123",
		StreamAccessRoleIds: []string{"user123", "user456"},
	}

	mockRepo.On("GetSetting").Return(testSetting, nil)

	result, err := useCase.GetSetting(ctx, role.Admin)

	assert.Equal(t, testSetting, result)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSystemSettingUseCase_GetSetting_Editor_Unauthorized(t *testing.T) {
	_, _, useCase := setup()
	ctx := context.Background()

	result, err := useCase.GetSetting(ctx, role.Editor)

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestSystemSettingUseCase_GetSetting_UnauthorizedUser(t *testing.T) {
	mockRepo, _, useCase := setup()
	ctx := context.Background()

	mockRepo.On("GetSetting").Return(nil, errors.New("unauthorized"))

	result, err := useCase.GetSetting(ctx, role.User)

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestSystemSettingUseCase_GetSetting_Guest_Unauthorized(t *testing.T) {
	_, _, useCase := setup()
	ctx := context.Background()

	result, err := useCase.GetSetting(ctx, role.Guest)

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestSystemSettingUseCase_GetSetting_Anonymous_Unauthorized(t *testing.T) {
	_, _, useCase := setup()
	ctx := context.Background()

	result, err := useCase.GetSetting(ctx, role.Anonymous)

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestSystemSettingUseCase_SetSetting_AdminUser(t *testing.T) {
	mockRepo, _, useCase := setup()
	ctx := context.Background()

	testSetting := &system.Setting{
		EditorRoleId:        "editor123",
		StreamAccessRoleIds: []string{"user123", "user456"},
	}

	mockRepo.On("SetSetting", testSetting).Return(nil)

	err := useCase.SetSetting(ctx, testSetting, role.Admin)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSystemSettingUseCase_SetSetting_Editor_Unauthorized(t *testing.T) {
	_, _, useCase := setup()
	ctx := context.Background()

	testSetting := &system.Setting{
		EditorRoleId:        "editor123",
		StreamAccessRoleIds: []string{"user123", "user456"},
	}

	err := useCase.SetSetting(ctx, testSetting, role.Editor)

	assert.Error(t, err)
}

func TestSystemSettingUseCase_SetSetting_UnauthorizedUser(t *testing.T) {
	mockRepo, _, useCase := setup()
	ctx := context.Background()

	testSetting := &system.Setting{
		EditorRoleId:        "editor123",
		StreamAccessRoleIds: []string{"user123", "user456"},
	}

	mockRepo.On("SetSetting", testSetting).Return(errors.New("unauthorized"))

	err := useCase.SetSetting(ctx, testSetting, role.User)

	assert.Error(t, err)
}

func TestSystemSettingUseCase_SetSetting_Guest_Unauthorized(t *testing.T) {
	_, _, useCase := setup()
	ctx := context.Background()

	testSetting := &system.Setting{
		EditorRoleId:        "editor123",
		StreamAccessRoleIds: []string{"user123", "user456"},
	}

	err := useCase.SetSetting(ctx, testSetting, role.Guest)

	assert.Error(t, err)
}

func TestSystemSettingUseCase_SetSetting_Anonymous_Unauthorized(t *testing.T) {
	_, _, useCase := setup()
	ctx := context.Background()

	testSetting := &system.Setting{
		EditorRoleId:        "editor123",
		StreamAccessRoleIds: []string{"user123", "user456"},
	}

	err := useCase.SetSetting(ctx, testSetting, role.Anonymous)

	assert.Error(t, err)
}
