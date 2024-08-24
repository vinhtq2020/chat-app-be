package domain

import (
	"context"
)

type FriendService interface {
	Create(ctx context.Context, userId string, friendId string, relation string) (int64, error)
	Response(ctx context.Context, userId string, friendId string, action string) (int64, error)
	Cancel(ctx context.Context, userId string, friendId string) (int64, error)
	Unfriend(ctx context.Context, userId string, friendId string) (int64, error)
}
