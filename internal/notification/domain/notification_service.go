package domain

import "context"

type NotificationService interface {
	Notify(ctx context.Context, generateId func() string, userId string, content string, subscribersId []string) (int64, error)
}
