package autocomplete

import (
	"context"
	domain_aggregator "go-service/internal/autocomplete/aggregator/domain"
	"go-service/internal/autocomplete/worker/domain"
	"go-service/pkg/logger"

	"github.com/redis/go-redis/v9"
)

func AggeratedData(ctx context.Context, aggregatorService domain_aggregator.AggregatorService, workerService domain.WorkerService, rdb *redis.Client, logger *logger.Logger) {

	res, err := aggregatorService.AggregatedData(ctx)
	if err != nil {
		logger.LogError(err.Error(), nil)
		return
	}

	if res == 0 {
		logger.LogInfo("no search queries to aggregate", nil)
		return
	}

	data, err := aggregatorService.All(ctx)
	if err != nil {
		logger.LogError(err.Error(), nil)
		return
	}

	_, err = workerService.CreateTries(ctx, data)
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
}
