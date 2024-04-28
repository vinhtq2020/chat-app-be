package app

import (
	"context"
	"go-service/internal/auth"
	auth_domain "go-service/internal/auth/domain"
	"go-service/internal/autocomplete/aggregator"
	querysearch "go-service/internal/autocomplete/query_search"
	query_search_domain "go-service/internal/autocomplete/query_search/domain"
	"go-service/internal/autocomplete/worker"
	"go-service/internal/configs"
	"go-service/internal/room"
	room_domain "go-service/internal/room/domain"
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
	Auth        auth_domain.AuthTransport
	User        user_domain.UserTransport
	Room        room_domain.RoomTransport
	QuerySearch query_search_domain.QuerySearchTransport
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
	}
	sequenceService := sequence.NewSequenceService(db)

	auth := auth.NewAuthTransport(db, validate, logger, configs.AccessTokenSecretKey, toArray)
	room := room.NewRoomTransport(db, sequenceService.Next, upgrader, toArray)
	user := user.NewUserTransport(db, toArray)
	querySearch := querysearch.NewQuerySearch(logger, rdb)

	aggregatorService := aggregator.NewAggregatorService(mongoDB, "querySearch", logger)
	workerService := worker.NewWorkerService(mongoDB, logger)
	scheduler, err := cron.NewSchedule(1 * time.Minute)
	if err != nil {
		return nil, err
	}

	logCron := cron.NewCron()
	logCron.AddJob(scheduler, cron.JobFunc(func() {
		res, err := aggregatorService.AggregatedData(context.Background())
		if err != nil {
			logger.LogError(err.Error(), nil)
			return
		}

		if res == 0 {
			logger.LogInfo("no search queries to aggregate", nil)
			return
		}

		data, err := aggregatorService.All(context.Background())
		if err != nil {
			logger.LogError(err.Error(), nil)
			return
		}

		_, err = workerService.CreateTries(context.Background(), data)
		if err != nil {
			logger.LogError(err.Error(), nil)
			return
		}

		trie, err := workerService.LoadTries(ctx)
		if err != nil {
			logger.LogError(err.Error(), nil)
			return
		}
		// cache new trie on redis
		err = rdb.Set(ctx, "autocomplete-trie", trie, 0).Err()
		if err != nil {
			logger.LogError(err.Error(), nil)
			return
		}
		logger.LogInfo("aggregated autocomplete trie success", nil)
	}))
	go logCron.Start()

	return &App{
		Auth:        auth,
		User:        user,
		Room:        room,
		QuerySearch: querySearch,
	}, nil
}
