package domain

import (
	"context"
	"go-service/internal/autocomplete/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkerRepository interface {
	Load(ctx context.Context, filter bson.M) (*model.Trie, error)
	Insert(ctx context.Context, trie model.Trie) (int64, error)
	Delete(ctx context.Context, filter bson.M) (int64, error)
	InTransaction(ctx context.Context, callback func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error)
}
