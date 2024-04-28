package repository

import (
	"context"
	"go-service/internal/autocomplete/model"
	"go-service/pkg/convert"
	"go-service/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type AggregatorRepository struct {
	mongoDB    *mongo.Database
	collection *mongo.Collection
	logger     *logger.Logger
}

func NewLogAggregatorRepository(mongoDB *mongo.Database, collectionName string, logger *logger.Logger) *AggregatorRepository {
	collection := mongoDB.Collection(collectionName)
	return &AggregatorRepository{
		mongoDB:    mongoDB,
		collection: collection,
		logger:     logger,
	}
}

func (r *AggregatorRepository) All(ctx context.Context, filter interface{}) ([]model.QueryCount, error) {
	var res []model.QueryCount
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}
	err = cursor.All(context.TODO(), &res)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}
	return res, nil
}

func (r *AggregatorRepository) InsertMany(ctx context.Context, queryCount []interface{}) (int64, error) {
	_, err := r.collection.InsertMany(ctx, convert.ConvertArrayToInterfaceArray(queryCount))
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}

	return 1, nil
}

func (r *AggregatorRepository) Delete(ctx context.Context, filter interface{}) (int64, error) {
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return 1, nil
}
