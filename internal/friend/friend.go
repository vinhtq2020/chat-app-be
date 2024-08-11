package friend

import (
	"go-service/internal/friend/delivery"
	friend_domain "go-service/internal/friend/domain"
	"go-service/internal/friend/repository"
	"go-service/internal/friend/usecase"
	"go-service/internal/notification/domain"
	user_domain "go-service/internal/user/domain"
	"go-service/pkg/database/postgres"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

func NewFriendHandler(db *gorm.DB, userRepository user_domain.UserRepository, notificationService domain.NotificationService, logger *logger.Logger, buildParan func(int) string) friend_domain.FriendTransport {
	friendRqRepo := repository.NewRequestFriendRepository(db, "friend_requests", "users", logger, buildParan)
	friendRepo := repository.NewFriendRepository(db, logger, postgres.BuildParam)
	sv := usecase.NewFriendService(friendRepo, friendRqRepo, userRepository, notificationService, logger)
	handler := delivery.NewFriendHandler(sv, logger)
	return handler
}
