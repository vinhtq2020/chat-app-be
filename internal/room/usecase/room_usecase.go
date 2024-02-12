package usecase

import (
	"fmt"
	"go-service/internal/room/domain"
	"go-service/pkg/convert"

	"github.com/gin-gonic/gin"
)

type RoomUsecase struct {
	repository     domain.RoomRepository
	NextSequenceNo func(c *gin.Context, module string) (int64, error)
}

func NewRoomUsecase(repository domain.RoomRepository,
	NextSequenceNo func(c *gin.Context, module string) (int64, error),
) *RoomUsecase {
	return &RoomUsecase{
		repository:     repository,
		NextSequenceNo: NextSequenceNo,
	}
}

func getSequenceNo(sequence int64) string {
	return fmt.Sprintf("room-%v", sequence)
}

func (u *RoomUsecase) All(c *gin.Context) ([]domain.Room, error) {
	return u.repository.All(c)
}

func (u *RoomUsecase) Load(c *gin.Context, id string) (*domain.Room, error) {
	room, err := u.repository.Load(c, id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (u *RoomUsecase) Create(c *gin.Context, room domain.Room) (int64, error) {
	sequence, err := u.NextSequenceNo(c, "room")
	if err != nil {
		return -1, err
	}

	room.Id = getSequenceNo(sequence)
	res, err := u.repository.Create(c, room)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (u *RoomUsecase) Patch(c *gin.Context, room domain.Room) (int64, error) {
	currentRoom, err := u.repository.Load(c, room.Id)
	if err != nil || currentRoom == nil {
		return 0, err
	}
	mp := convert.ToMapOmitEmpty(room)
	return u.repository.Patch(c, mp)
}

func (u *RoomUsecase) Delete(c *gin.Context, id string) (int64, error) {
	return u.repository.Delete(c, id)
}
