package main

import (
	"context"
	"fmt"
	"go-service/app"
	route "go-service/internal"
	"go-service/internal/configs"
	"go-service/pkg/database/mongo"
	"go-service/pkg/handler_fnc"
	"go-service/pkg/logger"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

func main() {
	mux := http.NewServeMux()
	ctx := context.Background()
	logger := logger.NewLogger()

	configs, err := getConfig(logger)
	if err != nil {
		logger.LogError(err.Error(), nil)
		return
	}
	ctx = context.WithValue(ctx, "secretKey", configs.AccessTokenSecretKey)

	mongoClient, err := mongo.NewMongoClient(context.Background(), mongo.MongoConfig{URI: configs.MongoConfig.URI})
	if err != nil {
		logger.LogError(err.Error(), nil)
		return
	}

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.RedisConfig.Addr,
		Password: configs.RedisConfig.Password,
		DB:       configs.RedisConfig.DB,
	})

	err = rdb.Ping(ctx).Err()
	if err != nil {
		logger.LogError(err.Error(), nil)
		return
	} else {
		fmt.Println("Pinged your deployment. You successly connected to Redis!")
	}

	defer func() {
		if err := rdb.Close(); err != nil {
			panic(err)
		}
	}()

	app, err := app.NewApp(ctx, mongoClient, rdb, configs, logger)
	if err != nil {
		logger.LogError(err.Error(), nil)
		return
	}

	route.Route(ctx, mux, app, logger, app.Auth.RefreshToken)
	err = http.ListenAndServe(":8080", handler_fnc.LogRequestHandler(mux, logger))
	if err != nil {
		logger.LogError(err.Error(), nil)
		return
	}

}

func getConfig(logger *logger.Logger) (configs.Config, error) {
	configs := configs.Config{}
	configFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return configs, err
	}

	err = yaml.Unmarshal(configFile, &configs)
	if err != nil {
		logger.LogError(err.Error(), nil)
		return configs, err
	}
	return configs, nil
}
