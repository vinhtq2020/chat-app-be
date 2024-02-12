package domain

import (
	"github.com/gin-gonic/gin"
)

type RoomService interface {
	All(c *gin.Context) ([]Room, error)
	Load(c *gin.Context, id string) (*Room, error)
	Create(c *gin.Context, room Room) (int64, error)
	Delete(c *gin.Context, id string) (int64, error)
	Patch(c *gin.Context, room Room) (int64, error)
}
