package domain

import (
	"context"
)

type NotificationRepository interface {
	Search(ctx context.Context, filter NotificationFilter) ([]Notification, error)
	Load(ctx context.Context, notificationId string) (*Notification, error)
	Insert(ctx context.Context, notification Notification) (int64, error)
	Patch(ctx context.Context, notification map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Total(ctx context.Context, clientID string) (int64, error)
	TotalUnread(ctx context.Context, clientID string) (int64, error)
}
