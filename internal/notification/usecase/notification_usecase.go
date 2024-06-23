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
	res, err := n.storageRepository.Insert(ctx, domain.Notification{
		Id:          generateId(),
		RequestorId: requestorId,
		Subscribers: subscriberIds,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Content:     content,
	})
	if err != nil {
		return res, err
	}

	for _, v := range subscriberIds {

	}
}
