package domain

import (
	"context"
)

type FriendService interface {
	SendFriendRequest(ctx context.Context, friendRequest FriendRequest) (int64, error)
	Patch(ctx context.Context, friendRequest map[string]interface{}) (int64, error)
	Delete(ctx context.Context, friendId string) (int64, error)
}
