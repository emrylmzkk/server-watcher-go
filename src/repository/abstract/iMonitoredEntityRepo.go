package repositoryAbstarct

import (
	"server-watcher-app/src/models"
	enumModels "server-watcher-app/src/models/enum"
	"time"
)

type MonitoredEntityRepository interface {
	CreateOrUpdate(entity *models.MonitoredEntity) error
	GetAllByType(pType enumModels.ProcessType) ([]models.MonitoredEntity, error)
	GetByExternalID(id string) (*models.MonitoredEntity, error)
	DeleteOldRecords(threshold time.Time) error
}
