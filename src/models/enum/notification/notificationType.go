package enumNotification

type NotificationType string

const (
	TypeCPUUsage  NotificationType = "CPU_USAGE"
	TypeDiskUsage NotificationType = "DISK_USAGE"
	TypeSystem    NotificationType = "SYSTEM_ALERTS"
)
