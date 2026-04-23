package models

import (
	enumModels "server-watcher-app/src/models/enum"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string              `gorm:"uniqueIndex;not null"`
	Password string              `gorm:"not null"`
	Name     string              `gorm:"not null"`
	Surname  string              `gorm:"not null"`
	UserRole enumModels.UserRole `gorm:"not null; default:2"`
}
