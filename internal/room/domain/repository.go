package domain

import (
	"context"
)

type RoomRepository interface {
	Load(ctx context.Context, id string) (*Room, error)
	All(ctx context.Context) ([]Room, error)
	Create(ctx context.Context, Room Room) (int64, error)
	Update(ctx context.Context, Room Room)
	Patch(ctx context.Context, Room map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
