package modelsDTOs

type ServerStatusSettingDTO struct {
	CPUThreshold         float64 `json:"cpu_threshold" validate:"required"`
	RAMThreshold         float64 `json:"ram_threshold" validate:"required"`
	DISKThreshold        float64 `json:"disk_threshold" validate:"required"`
	NotificationCooldown int     `json:"cooldown_minute" validate:"required"`
	IsDefault            bool    `json:"is_default" default:"false"`
}
