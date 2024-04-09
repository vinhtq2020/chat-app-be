package auth

import (
	"go-service/internal/auth/delivery"
	"go-service/internal/auth/domain"
	"go-service/internal/auth/repository"
	"go-service/internal/auth/usecase"
	"go-service/internal/auth/validator"
	user_repo "go-service/internal/user/repository"
	"go-service/pkg/database/sql/pq"
	"go-service/pkg/validate"

	"gorm.io/gorm"
)

func NewAuthTransport(db *gorm.DB, validate validate.Validate, secretKey string, toArray pq.Array) domain.AuthTransport {
	rp := repository.NewAuthRepository(db, "users_login_data", toArray)
	authValidator := validator.NewAuthValidator(db, "users_login_data", validate, toArray)
	userInfoRepository := user_repo.NewUserRepository(db, "users")
	sv := usecase.NewAuthUsecase(rp, authValidator, userInfoRepository, secretKey)
	hl := delivery.NewAuthHandler(sv)
	return hl
}
