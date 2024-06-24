package domain

import "context"

type NotificationStorageRepository interface {
	Total(ctx context.Context, clientID string) (int64, error)
	TotalUnread(ctx context.Context, clientID string) (int64, error)
	Insert(ctx context.Context, notification Notification) (int64, error)
	Patch(ctx context.Context, notification map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
