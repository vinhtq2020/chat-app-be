package usecase

import (
	"context"
	"fmt"
	"go-service/internal/room/domain"
	"go-service/pkg/convert"
)

type RoomUsecase struct {
	repository     domain.RoomRepository
	NextSequenceNo func(ctx context.Context, module string) (int64, error)
}

func NewRoomUsecase(repository domain.RoomRepository,
	NextSequenceNo func(ctx context.Context, module string) (int64, error),
) *RoomUsecase {
	return &RoomUsecase{
		repository:     repository,
		NextSequenceNo: NextSequenceNo,
	}
}

func getSequenceNo(sequence int64) string {
	return fmt.Sprintf("room-%v", sequence)
}

func (u *RoomUsecase) All(ctx context.Context) ([]domain.Room, error) {
	return u.repository.All(ctx)
}

func (u *RoomUsecase) Load(ctx context.Context, id string) (*domain.Room, error) {
	room, err := u.repository.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (u *RoomUsecase) Create(ctx context.Context, room domain.Room) (int64, error) {
	sequence, err := u.NextSequenceNo(ctx, "room")
	if err != nil {
		return -1, err
	}

	room.Id = getSequenceNo(sequence)
	res, err := u.repository.Create(ctx, room)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (u *RoomUsecase) Patch(ctx context.Context, room domain.Room) (int64, error) {
	currentRoom, err := u.repository.Load(ctx, room.Id)
	if err != nil || currentRoom == nil {
		return 0, err
	}
	mp := convert.ToMapOmitEmpty(room)
	return u.repository.Patch(ctx, mp)
}

func (u *RoomUsecase) Delete(ctx context.Context, id string) (int64, error) {
	return u.repository.Delete(ctx, id)
}
