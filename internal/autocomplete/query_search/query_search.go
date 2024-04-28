package querysearch

import (
	"go-service/internal/autocomplete/query_search/delivery"
	"go-service/internal/autocomplete/query_search/domain"
	"go-service/internal/autocomplete/query_search/repository"
	"go-service/internal/autocomplete/query_search/usecase"
	"go-service/pkg/logger"

	"github.com/redis/go-redis/v9"
)

func NewQuerySearch(logger *logger.Logger, rdb *redis.Client) domain.QuerySearchTransport {
	repository := repository.NewQuerySearchRepository(rdb, *logger)
	service := usecase.NewQuerySearchUsecase(repository)
	querySearch := delivery.NewQuerySearchHandler(logger, service)
	return querySearch
}
