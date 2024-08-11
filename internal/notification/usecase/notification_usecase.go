package usecase

import (
	"context"
	"go-service/internal/notification/domain"
	domain_notification "go-service/internal/notification/domain"
	"go-service/pkg/convert"
	"time"
)

type notificationService struct {
	broastcast        chan domain.Message
	storageRepository domain_notification.NotificationRepository
}

func NewNotificationService(storageRepository domain_notification.NotificationRepository, broastcast chan domain.Message) *notificationService {
	return &notificationService{
		storageRepository: storageRepository,
		broastcast:        broastcast,
	}
}

func (sv *notificationService) UpdateNotify(ctx context.Context, notificationId string, title string, content string, updatedBy string, notiType string, url *string, readed bool) (int64, error) {
	notification, err := sv.storageRepository.Load(ctx, notificationId)
	if err != nil || notification == nil {
		return 0, err
	}
	notification.Title = title
	notification.Content = content
	notification.Type = notiType
	notification.URL = url
	notification.UpdatedBy = updatedBy
	notification.UpdatedAt = time.Now()

	for i := range notification.Subscribers {
		notification.Subscribers[i].Readed = readed
	}
	params := convert.ToMapOmitEmpty(notification)

	res, err := sv.storageRepository.Patch(ctx, params)
	if err != nil || res <= 0 {
		return res, err
	}

	sv.broastcast <- domain.Message{
		Name: domain_notification.UPDATED,
		Data: *notification,
	}

	return res, nil
}

func (sv *notificationService) Search(ctx context.Context, filter domain.NotificationFilter) ([]domain.Notification, error) {
	return sv.storageRepository.Search(ctx, filter)
}

func (sv *notificationService) Notify(ctx context.Context, generateId func() string, requester domain.Requester, title string, content string, subscriberIds []string, url *string, notificationType string) (int64, error) {
	var subscribers []domain_notification.Subscriber
	for _, v := range subscriberIds {
		subscribers = append(subscribers, domain_notification.Subscriber{
			Id:     v,
			Readed: false,
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
		Deleted:     true,
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

func (sv *notificationService) Delete(ctx context.Context, id string) (int64, error) {
	res, err := sv.storageRepository.Delete(ctx, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (sv *notificationService) Patch(ctx context.Context, notificationId string, title string, content string, userId string) (int64, error) {
	oldNotification, err := sv.storageRepository.Load(ctx, notificationId)
	if err != nil {
		return -1, err
	}
	oldNotification.UpdatedAt = time.Now()
	oldNotification.UpdatedBy = userId
	oldNotification.Title = title
	oldNotification.Content = content

	newNotification := convert.ToMapOmitEmpty(oldNotification)
	res, err := sv.storageRepository.Patch(ctx, newNotification)
	if err != nil {
		return res, err
	}
	return res, nil
}
