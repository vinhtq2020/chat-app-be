package notification

import (
	"go-service/internal/notification/delivery"
	"go-service/internal/notification/domain"
	domain_notification "go-service/internal/notification/domain"
	"go-service/internal/notification/repository"
	"go-service/internal/notification/usecase"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func NewNotificationHandler(upgrader websocket.Upgrader, service domain.NotificationService, broastcast chan domain.Message, logger *logger.Logger) domain_notification.NotificacationTransport {
	handler := delivery.NewNotificationHandler(&upgrader, service, broastcast, logger)
	return handler
}

func NewNotificationService(db *gorm.DB, broastcast chan domain.Message, logger *logger.Logger, buildParam func(int) string, toArray pq.Array) domain_notification.NotificationService {
	notificationRepository := repository.NewNotificationRepository("notifications", buildParam, db, logger, toArray)
	notificationService := usecase.NewNotificationService(notificationRepository, broastcast)
	return notificationService
}
