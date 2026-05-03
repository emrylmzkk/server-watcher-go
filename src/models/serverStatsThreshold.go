package models

import "gorm.io/gorm"

type SSTSettings struct {
	gorm.Model
	UserID               uint    `gorm:"uniqueIndex"`
	CPUThreshold         float64 `gorm:"default:80.0"`
	RAMThreshold         float64 `gorm:"default:80.0"`
	DISKThreshold        float64 `gorm:"default:80.0"`
	NotificationCooldown int     `gorm:"default:5"`
}
