package friend

import (
	"go-service/internal/friend/delivery"
	"go-service/internal/friend/domain"
	"go-service/internal/friend/repository"
	"go-service/internal/friend/usecase"
	"go-service/pkg/database/postgres"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

func NewFriendHandler(db *gorm.DB, logger logger.Logger) domain.FriendTransport {
	friendRqRepo := repository.NewRequestFriendRepository(db, "friend_request", logger, postgres.BuildParam)
	sv := usecase.NewFriendService(friendRqRepo)
	handler := delivery.NewFriendHandler(sv)
	return handler
}
