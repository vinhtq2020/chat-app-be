package usecase

import (
	"context"
	domain_notification "go-service/internal/notification/notification_domain"
	"time"
)

type notificationService struct {
	storageRepository domain_notification.NotificationRepository
}

func NewNotificationService(storageRepository domain_notification.NotificationRepository) domain_notification.NotificationService {
	return &notificationService{
		storageRepository: storageRepository,
	}
}

func (n *notificationService) Notify(ctx context.Context, generateId func() string, requestorId string, content string, subscriberIds []string) (int64, error) {
	var subscribers []domain_notification.Subscriber
	for _, v := range subscriberIds {
		subscribers = append(subscribers, domain_notification.Subscriber{
			Id:     v,
			IsRead: false,
		})

	}

	res, err := n.storageRepository.Insert(ctx, domain_notification.Notification{
		Id:          generateId(),
		RequestorId: requestorId,
		Subscribers: subscribers,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     content,
	})

	if err != nil {
		return res, err
	}
	panic("")
}
