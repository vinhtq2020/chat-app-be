package room

import (
	"go-service/internal/room/delivery"
	"go-service/internal/room/domain"
	"go-service/internal/room/repository"
	"go-service/internal/room/usecase"
	"go-service/pkg/sql/pq"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func NewRoomTransport(db *gorm.DB, NextSequenceNo func(c *gin.Context, module string) (int64, error), upgrader websocket.Upgrader, toArray pq.Array) domain.RoomTransport {
	repo := repository.NewRoomReposiory(db, "rooms", toArray)
	service := usecase.NewRoomUsecase(repo, NextSequenceNo)
	handler := delivery.NewRoomHandler(service, upgrader)
	return handler
}
