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

func (sv *notificationService) Search(ctx context.Context, filter domain.NotificationFilter) ([]domain.Notification, error) {
	return sv.storageRepository.Search(ctx, filter)
}

func (sv *notificationService) Notify(ctx context.Context, generateId func() string, requester domain.Requester, title string, content string, subscriberIds []string, url *string, notificationType string) (int64, error) {
	var subscribers []domain_notification.Subscriber
	for _, v := range subscriberIds {
		subscribers = append(subscribers, domain_notification.Subscriber{
			Id:     v,
			IsRead: false,
		})

	}

	notification := domain_notification.Notification{
		Id:          generateId(),
		Requester:   requester,
		Title:       title,
		Content:     content,
		Subscribers: subscribers,
		Type:        notificationType,
		URL:         url,
		CreatedAt:   time.Now(),
		CreatedBy:   requester.Id,
		UpdatedAt:   time.Now(),
		UpdatedBy:   requester.Id,
		Visible:     true,
	}

	res, err := sv.storageRepository.Insert(ctx, notification)
	if err != nil {
		return -1, err
	}

	sv.broastcast <- domain.Message{
		Name: "notified",
		Data: notification,
	}
	return res, nil

}
