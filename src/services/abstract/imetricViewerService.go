package servicesAbstarct

import (
	"context"
	modelsDTOs "server-watcher-app/src/models/dtos"
)

type MetricViewerService interface {
	GetLogsFromDB(ctx context.Context, dto *modelsDTOs.ContainerLogRequestDTO) (*modelsDTOs.LogResponseDTO, error)
}
