package domain

import "context"

type NotificationService interface {
	Load(ctx context.Context, id string) (*Notification, error)
	Notify(ctx context.Context, generateId func() string, requester Requester, title string, content string, subscribersId []string, url *string, notificationType string) (int64, error)
	Search(ctx context.Context, filter NotificationFilter) ([]Notification, error)
	UpdateNotify(ctx context.Context, notificationId string, title string, content string, updatedBy string, notiType string, url *string, readed bool) (int64, error)
	Delete(ctx context.Context, notificationId string) (int64, error)
}
