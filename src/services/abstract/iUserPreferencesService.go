package servicesAbstarct

import (
	"context"
	modelsDTOs "server-watcher-app/src/models/dtos"
)

type IUserPreferencesService interface {
	CreateServerStatPreference(ctx context.Context, userID uint, dto *modelsDTOs.ServerStatusSettingDTO) (bool, error)
	ClearServerStatPreference(ctx context.Context, userID uint) (bool, error)
	GetUserSSTSettings(ctx context.Context, userID uint) (*modelsDTOs.ServerStatusSettingDTO, error)
}
