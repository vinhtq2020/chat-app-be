package repository

import (
	"context"

	"go-service/internal/autocomplete/model"
	mongodb "go-service/pkg/database/mongo"
	"go-service/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkerRepository struct {
	mongoDB    *mongo.Database
	collection *mongo.Collection
	logger     *logger.Logger
}

func NewWorkerRepository(mongoDB *mongo.Database, collectionName string, logger *logger.Logger) *WorkerRepository {
	collection := mongoDB.Collection(collectionName)
	return &WorkerRepository{
		collection: collection,
		mongoDB:    mongoDB,
	}
}

func (rp *WorkerRepository) Load(ctx context.Context, filter bson.M) (*model.Trie, error) {
	var trie *model.Trie
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"createdAt", -1}})

	err := rp.collection.FindOne(ctx, filter, findOptions).Decode(&trie)
	if err != nil {
		rp.logger.LogError(err.Error(), nil)
		return nil, err
	}
	return trie, nil
}

func (rp *WorkerRepository) Insert(ctx context.Context, trie model.Trie) (int64, error) {
	res, err := rp.collection.InsertOne(ctx, trie)
	if err != nil {
		rp.logger.LogError(err.Error(), nil)
		return -1, err
	}

	if res.InsertedID == nil {
		return 0, nil
	}
	return 1, nil
}

func (rp *WorkerRepository) Delete(ctx context.Context, filter bson.M) (int64, error) {
	res, err := rp.collection.DeleteMany(ctx, filter)
	if err != nil {
		rp.logger.LogError(err.Error(), nil)
		return -1, err
	}
	if res.DeletedCount <= 0 {
		return 0, nil
	}
	return res.DeletedCount, nil
}

func (rp *WorkerRepository) InTransaction(ctx context.Context, callback func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	return mongodb.InTransaction(ctx, rp.mongoDB.Client(), callback)
}
