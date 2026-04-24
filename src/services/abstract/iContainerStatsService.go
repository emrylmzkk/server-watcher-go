package servicesAbstarct

import (
	"context"
	"server-watcher-app/src/models"
	repositoryAbstarct "server-watcher-app/src/repository/abstract"
)

type ContainerStatsService interface {
	GetStatsByName(ctx context.Context, containerName string) ([]models.ContainerStatLog, error)
	GetStatsByNameP(ctx context.Context, p repositoryAbstarct.Pagination) (*repositoryAbstarct.PaginatedResult[models.ContainerStatLog], error)
}
