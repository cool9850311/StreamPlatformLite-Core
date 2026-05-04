package repository

import (
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/entity/system"
)

type SystemSettingRepository interface {
	GetSetting() (*system.Setting, error)
	SetSetting(setting *system.Setting) error
}
