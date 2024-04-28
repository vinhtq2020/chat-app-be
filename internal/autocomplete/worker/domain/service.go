package domain

import (
	"context"
	"go-service/internal/autocomplete/model"
)

type WorkerService interface {
	CreateTries(ctx context.Context, queryCounts []model.QueryCount) (int64, error)
	LoadTries(ctx context.Context) (*model.Trie, error)
}
