package auth

import (
	"go-service/internal/auth/delivery"
	"go-service/internal/auth/domain"
	"go-service/internal/auth/repository"
	"go-service/internal/auth/usecase"
	"go-service/internal/auth/validator"
	user_repo "go-service/internal/user/repository"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"
	"go-service/pkg/validate"

	"gorm.io/gorm"
)

func NewAuthTransport(db *gorm.DB, validate validate.Validate, logger *logger.Logger, secretKey string, toArray pq.Array) domain.AuthTransport {
	accountRepository := repository.NewAccountRepository(db, "users_login_data", logger, toArray)
	authValidator := validator.NewAuthValidator(db, "users_login_data", validate, toArray)
	userRepository := user_repo.NewUserRepository(db, "users", logger)
	refreshRepository := repository.NewRefreshTokenRepository(db, "refresh_tokens", logger, toArray)
	sv := usecase.NewAuthService(accountRepository, authValidator, userRepository, refreshRepository, secretKey)
	hl := delivery.NewAuthHandler(sv, logger)
	return hl
}
