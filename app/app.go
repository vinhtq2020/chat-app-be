package app

import (
	"go-service/internal/auth"
	auth_domain "go-service/internal/auth/domain"
	"go-service/internal/configs"
	"go-service/internal/room"
	room_domain "go-service/internal/room/domain"
	"go-service/internal/sequence"
	user_domain "go-service/internal/user/domain"
	"go-service/pkg/validate"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type App struct {
	Auth auth_domain.AuthTransport
	User user_domain.UserTransport
	Room room_domain.RoomTransport
	db   *gorm.DB
}

func NewApp(db *gorm.DB) (*App, error) {
	toArray := pq.Array
	validator := validator.New(validator.WithRequiredStructEnabled())
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	validate := validate.NewValidate(validator)

	configs := configs.Config{}
	configFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configFile, configs)
	if err != nil {
		return nil, err
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	sequenceService := sequence.NewSequenceService(db)

	auth := auth.NewAuthTransport(db, validate, configs.AccessTokenSecretKey, toArray)
	room := room.NewRoomTransport(db, sequenceService.Next, upgrader, toArray)

	return &App{
		Room: room,
		db:   db,
		Auth: auth,
	}, nil
}
