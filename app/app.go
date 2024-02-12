package app

import (
	"go-service/internal/auth"
	auth_domain "go-service/internal/auth/domain"
	"go-service/internal/room"
	room_domain "go-service/internal/room/domain"
	"go-service/internal/sequence"
	user_domain "go-service/internal/user/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type App struct {
	Auth auth_domain.AuthTransport
	User user_domain.UserTransport
	Room room_domain.RoomTransport
	db   *gorm.DB
}

func NewApp(db *gorm.DB) App {
	toArray := pq.Array
	validate := validator.New(validator.WithRequiredStructEnabled())

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	sequenceService := sequence.NewSequenceService(db)

	auth := auth.NewAuthTransport(db, validate)
	room := room.NewRoomTransport(db, sequenceService.Next, upgrader, toArray)

	return App{
		Room: room,
		db:   db,
		Auth: auth,
	}
}
