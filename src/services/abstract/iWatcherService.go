package servicesAbstarct

import (
	"context"
	modelsDTOs "server-watcher-app/src/models/dtos"
)

type WatcherService interface {
	SyncAll(ctx context.Context) error
	ControlProcess(ctx context.Context, dto *modelsDTOs.ActionOnProject) (bool, error)
}
