package domain

import "context"

type FriendRepository interface {
	AlreadyFriended(ctx context.Context, userId1 string, userId2 string) (int64, error)
	Upsert(ctx context.Context, friend Friend) (int64, error)
	UpdateStatus(ctx context.Context, userId1 string, userId2 string, status string) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
