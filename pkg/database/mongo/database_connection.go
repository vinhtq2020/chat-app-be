package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI string
}

func NewMongoClient(ctx context.Context, config MongoConfig) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI).SetServerAPIOptions(serverAPI))
	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client, nil
}
