package domain

import (
	"context"
)

type FriendRequestRepository interface {
	Exist(ctx context.Context, userId string, friendId string) (bool, error)
	All(ctx context.Context, userId string) ([]FriendRequest, error)
	Create(ctx context.Context, friendRq FriendRequest) (int64, error)
	Patch(ctx context.Context, friendRq map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
