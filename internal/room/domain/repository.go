package domain

import (
	"github.com/gin-gonic/gin"
)

type RoomRepository interface {
	Load(ctx *gin.Context, id string) (*Room, error)
	All(ctx *gin.Context) ([]Room, error)
	Create(ctx *gin.Context, Room Room) (int64, error)
	Update(ctx *gin.Context, Room Room)
	Patch(ctx *gin.Context, Room map[string]interface{}) (int64, error)
	Delete(ctx *gin.Context, id string) (int64, error)
}
