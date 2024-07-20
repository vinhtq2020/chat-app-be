package friend

import (
	"go-service/internal/friend/delivery"
	friend_domain "go-service/internal/friend/domain"
	"go-service/internal/friend/repository"
	"go-service/internal/friend/usecase"
	"go-service/internal/notification/domain"
	"go-service/pkg/database/postgres"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

func NewFriendHandler(db *gorm.DB, notificationService domain.NotificationService, logger *logger.Logger) friend_domain.FriendTransport {
	friendRqRepo := repository.NewRequestFriendRepository(db, "friend_requests", "users", logger, postgres.BuildParam)
	sv := usecase.NewFriendService(friendRqRepo, notificationService, logger)
	handler := delivery.NewFriendHandler(sv)
	return handler
}
