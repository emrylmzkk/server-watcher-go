package servicesAbstarct

import (
	"context"
	"server-watcher-app/src/models"
)

type IServerGeneralService interface {
	GetCPUPercent(ctx context.Context) (float64, error)
	//GetRAMStats() (used, total, percent float64, err error)
	GetRamStats(ctx context.Context) (used, total, percent float64, err error)
	GetDiskStats(ctx context.Context) (used, total, percent float64, err error)
	GetSystemStats(ctx context.Context) (*models.SystemStats, error)
	//GetSystemStats() (*SystemStats, error)
	GetStatsForNotification(ctx context.Context) (*models.SystemStats, error)
}
