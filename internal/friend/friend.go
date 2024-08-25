package friend

import (
	"go-service/internal/friend/delivery"
	friend_domain "go-service/internal/friend/domain"
	"go-service/internal/friend/repository"
	"go-service/internal/friend/usecase"
	"go-service/internal/notification/domain"
	sequence "go-service/internal/sequence/domain"
	user_domain "go-service/internal/user/domain"
	"go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

func NewFriendHandler(db *gorm.DB, userRepository user_domain.UserRepository, notificationService domain.NotificationService, sequence sequence.SequenceService, logger *logger.Logger, buildParan func(int) string, toArray pq.Array) friend_domain.FriendTransport {
	friendRepo := repository.NewFriendRepository(db, logger, postgres.BuildParam, toArray)
	sv := usecase.NewFriendService(friendRepo, userRepository, notificationService, sequence, logger)
	handler := delivery.NewFriendHandler(sv, logger)
	return handler
}
