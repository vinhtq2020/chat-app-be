package notification

import (
	"go-service/internal/notification/delivery"
	domain_notification "go-service/internal/notification/domain"
	"go-service/internal/notification/repository"
	"go-service/internal/notification/usecase"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func NewNotificationHandler(upgrader websocket.Upgrader, hub domain_notification.Hub, logger *logger.Logger) domain_notification.NotificacationTransport {
	go hub.Run()
	handler := delivery.NewNotificationHandler(&upgrader, hub, logger)
	return handler
}

func NewNotificationService(db *gorm.DB, hub domain_notification.Hub, logger *logger.Logger, buildParam func(int) string, toArray pq.Array) domain_notification.NotificationService {
	notificationRepository := repository.NewNotificationRepository("notifications", buildParam, db, logger, toArray)
	notificationService := usecase.NewNotificationService(notificationRepository, hub)
	return notificationService
}
