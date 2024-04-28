package domain

import (
	"context"
)

type RoomService interface {
	All(ctx context.Context) ([]Room, error)
	Load(ctx context.Context, id string) (*Room, error)
	Create(ctx context.Context, room Room) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Patch(ctx context.Context, room Room) (int64, error)
}
