package servicesAbstarct

import (
	"context"
	modelsDTOs "server-watcher-app/src/models/dtos"
)

type DockerService interface {
	GetContainers(ctx context.Context) ([]modelsDTOs.DockerContainerResponseDTO, error)
	StartContainer(ctx context.Context, containerID string) error
	StopContainer(ctx context.Context, containerID string) error
	StartStatsCollector(ctx context.Context)
	GetCachedStats() []modelsDTOs.DockerStats
	GetActiveContainers(ctx context.Context) ([]modelsDTOs.DockerContainerResponseDTO, error)
}
