package repositoryConcrete

import (
	"context"
	"server-watcher-app/src/models"
	enumModels "server-watcher-app/src/models/enum"
	repositoryAbstarct "server-watcher-app/src/repository/abstract"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	repositoryAbstarct.BaseRepository[models.MonitoredEntity]
	GetPm2Projects(ctx context.Context) ([]models.MonitoredEntity, error)
	GetPm2ByExternalId(ctx context.Context, externalId string) (models.MonitoredEntity, error)
}

type projectRepository struct {
	repositoryAbstarct.BaseRepository[models.MonitoredEntity]
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{
		BaseRepository: repositoryAbstarct.NewBaseRepository[models.MonitoredEntity](db),
	}
}

func (r *projectRepository) GetPm2Projects(ctx context.Context) ([]models.MonitoredEntity, error) {

	var projects []models.MonitoredEntity

	err := r.Query(ctx).Where("type = ?", enumModels.PM2).Find(&projects).Error

	if err != nil {
		return nil, err
	}

	return projects, nil

}

func (r *projectRepository) GetPm2ByExternalId(ctx context.Context, externalId string) (models.MonitoredEntity, error) {
	var project models.MonitoredEntity
	err := r.Query(ctx).
		Select("name", "project_start_command", "project_path", "external_id").
		Where("external_id = ?", externalId).
		Where("type = ?", enumModels.PM2).
		First(&project).Error
	return project, err
}

func (r *projectRepository) GetDockerProjects(ctx context.Context) ([]models.MonitoredEntity, error) {

	var projects []models.MonitoredEntity

	err := r.Query(ctx).Where("type = ?", enumModels.Docker).Find(&projects).Error

	if err != nil {
		return nil, err
	}

	return projects, nil

}
