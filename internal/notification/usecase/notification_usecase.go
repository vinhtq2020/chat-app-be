package notification

import (
	"context"
	"go-service/internal/notification/domain"
	"time"
)

type notificationService struct {
	storageRepository domain.NotificationStorageRepository
}

func NewNotificationService(storageRepository domain.NotificationStorageRepository) domain.NotificationService {
	return &notificationService{
		storageRepository: storageRepository,
	}
}

func (n *notificationService) Notify(ctx context.Context, generateId func() string, requestorId string, content string, subscriberIds []string) (int64, error) {
	var subscribers []domain.Subscriber
	for _, v := range subscriberIds {
		subscribers = append(subscribers, domain.Subscriber{
			SubscriberId: v,
			IsRead:       false,
		})

	}

	res, err := n.storageRepository.Insert(ctx, domain.Notification{
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
