package domain

import (
	"context"
)

type FriendService interface {
	SendFriendRequest(ctx context.Context, friendRequest FriendRequest) (int64, error)
	UpdateFriendRequest(ctx context.Context, requestId string, action string) (int64, error)
	Delete(ctx context.Context, friendId string) (int64, error)
}
