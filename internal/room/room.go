package room

import (
	"context"
	"go-service/internal/room/delivery"
	"go-service/internal/room/domain"
	"go-service/internal/room/repository"
	"go-service/internal/room/usecase"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func NewRoomTransport(db *gorm.DB, NextSequenceNo func(c context.Context, module string) (int64, error), upgrader websocket.Upgrader, logger *logger.Logger, toArray pq.Array) domain.RoomTransport {
	repo := repository.NewRoomReposiory(db, "rooms", toArray)
	service := usecase.NewRoomUsecase(repo, NextSequenceNo)
	handler := delivery.NewRoomHandler(service, upgrader, logger)
	return handler
}
