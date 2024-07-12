package user

import (
	"go-service/internal/user/delivery"
	"go-service/internal/user/domain"
	"go-service/internal/user/repository"
	"go-service/internal/user/usecase"
	"go-service/internal/utils/search"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

func NewUserTransport(db *gorm.DB, logger *logger.Logger, toArray pq.Array) domain.UserTransport {
	repo := repository.NewUserRepository(db, "users", logger)
	service := usecase.NewUserUsecase(repo)
	searchRepo := search.NewSearchRepository("users", db, toArray)
	searchService := search.NewSearchService[domain.User](searchRepo)
	handler := delivery.NewUserHandler(service, searchService)
	return handler
}
