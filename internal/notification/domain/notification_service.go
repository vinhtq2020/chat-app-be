package domain

import "context"

type NotificationService interface {
	Notify(ctx context.Context, generateId func() string, requester Requester, title string, content string, subscribersId []string, url *string, notificationType string) (int64, error)
	Search(ctx context.Context, filter NotificationFilter) ([]Notification, error)
}
