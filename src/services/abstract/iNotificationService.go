package servicesAbstarct

import enumNotification "server-watcher-app/src/models/enum/notification"

type INotificationService interface {
	SendToUser(userID uint, token string, title string, body string, cooldown int, notifType enumNotification.NotificationType) error
}
