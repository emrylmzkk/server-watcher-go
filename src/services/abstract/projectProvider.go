package servicesAbstarct

import (
	"context"
	"server-watcher-app/src/models"
	enumModels "server-watcher-app/src/models/enum"
)

type ProcessProvider interface {
	GetProviderType() enumModels.ProcessType
	ListProcesses(ctx context.Context) ([]models.MonitoredEntity, error)
	StopProcess(ctx context.Context, externalID string) error
	StartProcess(ctx context.Context, externalID string) error
	//GetLogs(ctx context.Context, externalID string) ([]string, error)
}
