package querysearch

import (
	"go-service/internal/query_search/delivery"
	"go-service/internal/query_search/domain"
	"go-service/pkg/logger"
)

func NewQuerySearch(logger *logger.Logger) domain.QuerySearchTransport {
	querySearch := delivery.NewQuerySearchHandler(logger)
	return querySearch
}
