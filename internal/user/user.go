package user

import (
	search_repository "go-service/internal/search/repository"
	search_service "go-service/internal/search/service"
	"go-service/internal/user/delivery"
	"go-service/internal/user/domain"
	"go-service/internal/user/repository"
	"go-service/internal/user/usecase"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

func NewUserTransport(db *gorm.DB, logger *logger.Logger, toArray pq.Array) domain.UserTransport {
	repo := repository.NewUserRepository(db, "users", logger)
	service := usecase.NewUserUsecase(repo)
	searchRepo := search_repository.NewSearchRepository("users", db, toArray)
	searchService := search_service.NewSearchService[domain.User](searchRepo)
	handler := delivery.NewUserHandler(service, searchService)
	return handler
}
