package generic

import (
	"server-watcher-app/src/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("wathcer_app.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Tabloları otomatik oluştur
	err = db.AutoMigrate(
		&models.MonitoredEntity{},
	)
	return db, err
}
