package usecase

import (
	"context"
	domain_notification "go-service/internal/notification/domain"
	"time"
)

type notificationService struct {
	hub               domain_notification.Hub
	storageRepository domain_notification.NotificationRepository
}

func NewNotificationService(storageRepository domain_notification.NotificationRepository, hub domain_notification.Hub) domain_notification.NotificationService {
	return &notificationService{
		storageRepository: storageRepository,
		hub:               hub,
	}
}

func (sv *notificationService) Notify(ctx context.Context, generateId func() string, requestorId string, content string, subscriberIds []string) (int64, error) {
	var subscribers []domain_notification.Subscriber
	for _, v := range subscriberIds {
		subscribers = append(subscribers, domain_notification.Subscriber{
			Id:     v,
			IsRead: false,
		})

	}

	notification := domain_notification.Notification{
		Id:          generateId(),
		RequestorId: requestorId,
		Subscribers: subscribers,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     content,
	}

	// res, err := sv.storageRepository.Insert(ctx, notification)
	// if err != nil {
	// 	return -1, err
	// }

	sv.hub.SendMessage(domain_notification.WsMsg{
		Name: "notified",
		Data: notification,
	})
	return 1, nil

}
