package notification

import (
	"go-service/internal/notification/delivery"
	domain_notification "go-service/internal/notification/notification_domain"
	"go-service/internal/notification/repository"
	"go-service/internal/notification/usecase"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func NewNotificationHandler(upgrader websocket.Upgrader, logger *logger.Logger) domain_notification.NotificacationTransport {
	handler := delivery.NewNotificationHandler(&upgrader, logger)
	return handler
}

func NewNotificationService(db *gorm.DB, buildParam func(int) string, toArray pq.Array) domain_notification.NotificationService {
	notificationRepository := repository.NewNotificationRepository("notifications", buildParam, db, toArray)
	notificationService := usecase.NewNotificationService(notificationRepository)
	return notificationService
}
