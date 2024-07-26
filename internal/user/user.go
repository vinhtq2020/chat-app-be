package user

import (
	"go-service/internal/user/delivery"
	user_domain "go-service/internal/user/domain"
	repository "go-service/internal/user/repository"
	"go-service/internal/user/usecase"
	"go-service/internal/utils/search"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

func NewUserTransport(db *gorm.DB, userRepository user_domain.UserRepository, toArray pq.Array) user_domain.UserTransport {
	service := usecase.NewUserUsecase(userRepository)
	searchRepo := search.NewSearchRepository("users", db, toArray)
	searchService := search.NewSearchService[user_domain.User](searchRepo)
	handler := delivery.NewUserHandler(service, searchService)
	return handler
}

func NewUserRepository(db *gorm.DB, logger *logger.Logger, toArray pq.Array) user_domain.UserRepository {
	return repository.NewUserRepository(db, "users", logger, toArray)
}
