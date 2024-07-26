package domain

import (
	"context"

	"gorm.io/gorm"
)

type FriendRequestRepository interface {
	Exist(ctx context.Context, userId string, friendId string) (bool, error)
	All(ctx context.Context, userId string) ([]FriendRequest, error)
	Total(ctx context.Context) (int64, error)
	Create(ctx context.Context, friendRq FriendRequest) (int64, error)
	Patch(ctx context.Context, friendRq map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	InTransaction(ctx context.Context, ex func(ctx context.Context, db *gorm.DB) (int64, error)) (int64, error)
}
