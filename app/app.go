package app

import (
	"context"
	"go-service/internal/auth"
	auth_domain "go-service/internal/auth/domain"
	"go-service/internal/autocomplete"
	"go-service/internal/autocomplete/aggregator"
	querysearch "go-service/internal/autocomplete/query_search"
	query_search_domain "go-service/internal/autocomplete/query_search/domain"
	"go-service/internal/autocomplete/worker"
	"go-service/internal/configs"
	"go-service/internal/friend"
	friend_domain "go-service/internal/friend/domain"
	"go-service/internal/notification"
	"go-service/internal/notification/domain"
	"net/http"

	"go-service/internal/room"
	room_domain "go-service/internal/room/domain"
	"go-service/internal/search_tool"
	"go-service/internal/sequence"
	"go-service/internal/user"
	user_domain "go-service/internal/user/domain"
	"go-service/pkg/cron"
	"go-service/pkg/database/postgres"
	"go-service/pkg/logger"
	"go-service/pkg/validate"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type serviceKey string

const (
	authServiceKey serviceKey = "authService"
)

type App struct {
	Auth         auth_domain.AuthTransport
	User         user_domain.UserTransport
	Room         room_domain.RoomTransport
	QuerySearch  query_search_domain.QuerySearchTransport
	SearchTool   search_tool.SearchToolTransport
	Notification domain.NotificacationTransport
	Friend       friend_domain.FriendTransport
}

func NewApp(ctx context.Context, mongoClient *mongo.Client, rdb *redis.Client, configs configs.Config, logger *logger.Logger) (*App, error) {
	db, err := postgres.NewPostgresDb(configs.PostgresConfig.DSN)
	if err != nil {
		return nil, err
	}

	mongoDB := mongoClient.Database(configs.MongoConfig.Database)
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

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	sequenceService := sequence.NewSequenceService(db)

	auth := auth.NewAuthTransport(db, validate, logger, configs.AccessTokenSecretKey, toArray)
	room := room.NewRoomTransport(db, sequenceService.Next, upgrader, logger, toArray)

	userRepository := user.NewUserRepository(db, logger, toArray)
	user := user.NewUserTransport(db, userRepository, toArray)
	querySearch := querysearch.NewQuerySearch(logger, rdb)

	aggregatorService := aggregator.NewAggregatorService(mongoDB, "querySearch", logger)
	workerService := worker.NewWorkerService(mongoDB, logger)
	scheduler, err := cron.NewSchedule(24 * 7 * time.Hour)
	if err != nil {
		return nil, err
	}

	logCron := cron.NewCron()
	logCron.AddJob(scheduler, cron.JobFunc(func() { autocomplete.AggeratedData(ctx, aggregatorService, workerService, rdb, logger) }))
	go logCron.Start()

	go room.HandleMessages()

	broastcast := make(chan domain.Message)
	notificationService := notification.NewNotificationService(db, broastcast, logger, postgres.BuildParam, toArray)
	notification := notification.NewNotificationHandler(upgrader, notificationService, broastcast, logger)
	go notification.HandleMessages()
	searchTool := search_tool.NewSearchToolTransport(db, postgres.BuildParam, logger, toArray)

	friend := friend.NewFriendHandler(db, userRepository, notificationService, logger, postgres.BuildParam)
	return &App{
		Auth:         auth,
		User:         user,
		Room:         room,
		QuerySearch:  querySearch,
		SearchTool:   searchTool,
		Notification: notification,
		Friend:       friend,
	}, nil
}
