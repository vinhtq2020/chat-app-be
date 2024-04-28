package repository

import (
	"context"
	"go-service/internal/autocomplete/model"
	"go-service/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type QuerySearchRepository struct {
	rdb    *redis.Client
	logger logger.Logger
}

func NewQuerySearchRepository(rdb *redis.Client, logger logger.Logger) *QuerySearchRepository {
	return &QuerySearchRepository{
		rdb:    rdb,
		logger: logger,
	}
}

func (r *QuerySearchRepository) LoadTrie(ctx context.Context) (*model.Trie, error) {
	trie := model.Trie{}
	err := r.rdb.Get(ctx, "autocomplete-trie").Scan(&trie)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}

	return &trie, nil
}
