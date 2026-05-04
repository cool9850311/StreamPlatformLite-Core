package usecase

import (
	"context"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/interface/repository"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/entity/errors"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/entity/system"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/interface/logger"
	"github.com/cool9850311/StreamPlatformLite-Core/pkg/role"
)

type SystemSettingUseCase struct {
	systemSettingRepo repository.SystemSettingRepository
	Log               logger.Logger
}

func NewSystemSettingUseCase(systemSettingRepo repository.SystemSettingRepository, log logger.Logger) *SystemSettingUseCase {
	return &SystemSettingUseCase{
		systemSettingRepo: systemSettingRepo,
		Log:               log,
	}
}

func (u *SystemSettingUseCase) CheckRole(userRole role.Role) error {
	if userRole != role.Admin {
		return errors.ErrUnauthorized
	}
	return nil
}

func (u *SystemSettingUseCase) GetSetting(ctx context.Context, userRole role.Role) (*system.Setting, error) {
	if err := u.CheckRole(userRole); err != nil {
		u.Log.Error(ctx, "Unauthorized access to GetSetting")
		return nil, err
	}
	setting, err := u.systemSettingRepo.GetSetting()
	if err != nil {
		u.Log.Error(ctx, "Error getting system setting")
		return nil, err
	}
	return setting, nil
}

// SetSetting updates the system setting
func (u *SystemSettingUseCase) SetSetting(ctx context.Context, setting *system.Setting, userRole role.Role) error {
	if err := u.CheckRole(userRole); err != nil {
		u.Log.Error(ctx, "Unauthorized access to SetSetting")
		return err
	}
	return u.systemSettingRepo.SetSetting(setting)
}
