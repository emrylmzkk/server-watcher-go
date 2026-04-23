package models

import (
	enumModels "server-watcher-app/src/models/enum"
	"time"

	"gorm.io/gorm"
)

type MonitoredEntity struct {
	gorm.Model

	ID                  int    `gorm:"primaryKey"`
	ExternalID          string `gorm:"uniqueIndex"`
	Name                string
	Type                enumModels.ProcessType `gorm:"index"`
	Status              string
	LastCheck           time.Time
	ProjectPath         *string
	ProjectStartCommand *string
	ProjectRuntimeType  *enumModels.ProjectRuntimeType
}
