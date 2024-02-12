package user

import (
	"go-service/internal/user/delivery"
	"go-service/internal/user/domain"
	"go-service/internal/user/repository"
	"go-service/internal/user/usecase"

	"gorm.io/gorm"
)

func NewUserTransport(db *gorm.DB) domain.UserTransport {
	repo := repository.NewUserRepository(db, "users")
	service := usecase.NewUserUsecase(repo)
	handler := delivery.NewUserHandler(service)
	return handler
}
