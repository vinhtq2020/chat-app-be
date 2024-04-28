package domain

import (
	"context"
	"go-service/internal/autocomplete/model"
)

type AggregatorRepository interface {
	All(ctx context.Context, filter interface{}) ([]model.QueryCount, error)
	InsertMany(ctx context.Context, queryCount []interface{}) (int64, error)
	Delete(ctx context.Context, filter interface{}) (int64, error)
}
