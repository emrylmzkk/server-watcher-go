package repositoryConcrete

import (
	"context"
	"errors"
	"server-watcher-app/src/models"
	repositoryAbstarct "server-watcher-app/src/repository/abstract"

	"gorm.io/gorm"
)

type ContainerStatLogRepo interface {
	repositoryAbstarct.BaseRepository[models.ContainerStatLog]
	GetStatsWithContainerName(ctx context.Context, containerName string) ([]models.ContainerStatLog, error)
}

type containerStatLogRepo struct {
	repositoryAbstarct.BaseRepository[models.ContainerStatLog]
}

func NewContainerStatLogRepo(db *gorm.DB) ContainerStatLogRepo {
	return &containerStatLogRepo{
		BaseRepository: repositoryAbstarct.NewBaseRepository[models.ContainerStatLog](db),
	}
}

func (r *containerStatLogRepo) GetStatsWithContainerName(ctx context.Context, containerName string) ([]models.ContainerStatLog, error) {

	var containerStats []models.ContainerStatLog

	err := r.Query(ctx).Where("container_name = ?", containerName).Find(&containerStats).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return containerStats, nil
		}

		return nil, err
	}

	return containerStats, nil

}
