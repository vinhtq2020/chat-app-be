package aggregator

import (
	"go-service/internal/autocomplete/aggregator/domain"
	"go-service/internal/autocomplete/aggregator/repository"
	"go-service/internal/autocomplete/aggregator/usecase"
	"go-service/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewAggregatorService(mongoDb *mongo.Database, collection string, logger *logger.Logger) domain.AggregatorService {
	repo := repository.NewLogAggregatorRepository(mongoDb, collection, logger)
	aggregator := usecase.NewLogAggregatorService(repo, logger)
	return aggregator
}
