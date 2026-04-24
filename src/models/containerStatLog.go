package models

import (
	"time"

	"gorm.io/gorm"
)

type ContainerStatLog struct {
	ID            uint           `gorm:"primarykey"`
	ContainerName string         `gorm:"index"`
	CPU           float64
	MemoryMB      float64
	NetworkRX     uint64
	NetworkTX     uint64
	CreatedAt     time.Time      `gorm:"index"`
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
