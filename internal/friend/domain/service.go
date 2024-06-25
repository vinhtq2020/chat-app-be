package domain

import "context"

type FriendService interface {
	SendFriendRequest(ctx context.Context, userId string, friendId string) (int64, error)
	Patch(ctx context.Context, friendId string, status FriendRequestStatus) (int64, error)
	Delete(ctx context.Context, friendId string) (int64, error)
}
