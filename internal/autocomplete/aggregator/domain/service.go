package domain

import (
	"context"
	"go-service/internal/autocomplete/model"
)

type AggregatorService interface {
	AggregatedData(context.Context) (int64, error)
	All(ctx context.Context) ([]model.QueryCount, error)
}
