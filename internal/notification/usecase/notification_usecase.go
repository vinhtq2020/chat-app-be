package usecase

import (
	"context"
	"go-service/internal/notification/domain"
	domain_notification "go-service/internal/notification/domain"
	"time"
)

type notificationService struct {
	broastcast        chan domain.Message
	storageRepository domain_notification.NotificationRepository
}

func NewNotificationService(storageRepository domain_notification.NotificationRepository, broastcast chan domain.Message) domain_notification.NotificationService {
	return &notificationService{
		storageRepository: storageRepository,
		broastcast:        broastcast,
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

	sv.broastcast <- domain.Message{
		Name: "notified",
		Data: notification,
	}
	return 1, nil

}
