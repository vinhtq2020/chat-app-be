package auth

import (
	"go-service/internal/auth/delivery"
	"go-service/internal/auth/domain"
	"go-service/internal/auth/repository"
	"go-service/internal/auth/usecase"
	"go-service/internal/auth/validator"
	user_repo "go-service/internal/user/repository"

	validate "github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func NewAuthTransport(db *gorm.DB, validate *validate.Validate) domain.AuthTransport {
	rp := repository.NewAuthRepository(db, "users_login_data")
	authValidator := validator.NewAuthValidator(validate)
	userInfoRepository := user_repo.NewUserRepository(db, "users")
	sv := usecase.NewAuthUsecase(rp, authValidator, userInfoRepository)
	hl := delivery.NewAuthHandler(sv)
	return hl
}
