package worker

import (
	"go-service/internal/autocomplete/worker/domain"
	"go-service/internal/autocomplete/worker/repository"
	"go-service/internal/autocomplete/worker/usecase"
	"go-service/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewWorkerService(db *mongo.Database, logger *logger.Logger) domain.WorkerService {
	repository := repository.NewWorkerRepository(db, "tries", logger)
	service := usecase.NewWorkerService(repository)
	return service
}
