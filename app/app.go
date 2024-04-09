package app

import (
	"context"
	"fmt"
	"go-service/internal/auth"
	auth_domain "go-service/internal/auth/domain"
	"go-service/internal/configs"
	querysearch "go-service/internal/query_search"
	query_search_domain "go-service/internal/query_search/domain"
	"go-service/internal/room"
	room_domain "go-service/internal/room/domain"
	"go-service/internal/sequence"
	"go-service/internal/user"
	user_domain "go-service/internal/user/domain"
	"go-service/pkg/cron"
	"go-service/pkg/database/mongo"
	"go-service/pkg/database/postgres"
	"go-service/pkg/logger"
	"go-service/pkg/validate"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type App struct {
	Auth        auth_domain.AuthTransport
	User        user_domain.UserTransport
	Room        room_domain.RoomTransport
	QuerySearch query_search_domain.QuerySearchTransport
}

func NewApp(logger *logger.Logger) (*App, error) {
	db, err := postgres.NewPostgresDb()
	if err != nil {
		return nil, err
	}

	_, err = mongo.NewMongoClient(context.Background(), mongo.MongoConfig{URI: "mongodb://root:Abcd1234@mongo:27017/"})
	if err != nil {
		return nil, err

	}
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
	user := user.NewUserTransport(db, toArray)
	querySearch := querysearch.NewQuerySearch(logger)

	scheduler, err := cron.NewSchedule(3 * time.Second)
	if err != nil {
		return nil, err
	}
	logCron := cron.NewCron()
	logCron.AddJob(scheduler, cron.JobFunc(func() {
		fmt.Println(time.Now())
	}))
	go logCron.Start()
	return &App{
		Auth:        auth,
		User:        user,
		Room:        room,
		QuerySearch: querySearch,
	}, nil
}
