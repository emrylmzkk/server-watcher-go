package repositoryConcrete

import (
	"server-watcher-app/src/models"
	enumModels "server-watcher-app/src/models/enum"
	repository "server-watcher-app/src/repository/abstract"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sqliteRepo struct {
	db *gorm.DB
}

func NewSqliteRepository(db *gorm.DB) repository.MonitoredEntityRepository {
	return &sqliteRepo{db: db}
}

func (r *sqliteRepo) CreateOrUpdate(entity *models.MonitoredEntity) error {
	// Upsert işlemi: Eğer ID varsa güncelle, yoksa oluştur
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "last_check", "name"}),
	}).Create(entity).Error
}

func (r *sqliteRepo) GetAllByType(pType enumModels.ProcessType) ([]models.MonitoredEntity, error) {
	var entities []models.MonitoredEntity
	err := r.db.Where("type = ?", pType).Find(&entities).Error
	return entities, err
}

func (r *sqliteRepo) GetByExternalID(id string) (*models.MonitoredEntity, error) {
	var entity models.MonitoredEntity
	err := r.db.Where("external_id = ?", id).First(&entity).Error
	return &entity, err
}

func (r *sqliteRepo) DeleteOldRecords(threshold time.Time) error {
	return r.db.Unscoped().Where("last_check < ?", threshold).Delete(&models.MonitoredEntity{}).Error
}
