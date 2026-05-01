package servicesConcrete

import (
	"context"
	"fmt"
	"log"
	enumNotification "server-watcher-app/src/models/enum/notification"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"sync"
	"time"

	"firebase.google.com/go/v4/messaging"
)

type notificationService struct {
	fcmClient *messaging.Client

	cooldowns sync.Map
}

func NewNotificationService(fcmClient *messaging.Client) servicesAbstarct.INotificationService {
	return &notificationService{
		fcmClient: fcmClient,
	}
}

func (s *notificationService) SendToUser(userID uint, token string, title string, body string, notifType enumNotification.NotificationType) error {

	if token == "" {
		log.Printf("[FCM] Uyarı: Kullanıcı %v için token boş, gönderim iptal edildi.", userID)
		return nil
	}

	cacheKey := fmt.Sprintf("%d_%s", userID, notifType)

	if lastSent, ok := s.cooldowns.Load(cacheKey); ok {
		if time.Since(lastSent.(time.Time)) < 5*time.Minute {
			log.Printf("[FCM] Spam Engellendi: Kullanıcı %d, Tip %s (Henüz 5 dk dolmadı)", userID, notifType)
			return nil
		}
	}

	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	log.Printf("[FCM] Bildirim gönderiliyor -> User: %d, Title: %s", userID, title)

	msgID, err := s.fcmClient.Send(context.Background(), message)

	if err != nil {
		log.Printf("[FCM] HATA: Bildirim gönderilemedi: %v", err)
		return nil
	}

	s.cooldowns.Store(cacheKey, time.Now())

	log.Printf("[FCM] Başarılı! Mesaj ID: %s", msgID)

	return nil

}
