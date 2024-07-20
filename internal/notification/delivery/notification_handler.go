package delivery

import (
	domain_notification "go-service/internal/notification/domain"
	"go-service/pkg/logger"
	"net/http"

	"github.com/gorilla/websocket"
)

type NotificationHandler struct {
	hub      domain_notification.Hub
	upgrader *websocket.Upgrader
	logger   *logger.Logger
}

func NewNotificationHandler(upgrader *websocket.Upgrader, hub domain_notification.Hub, logger *logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		upgrader: upgrader,
		logger:   logger,
		hub:      hub,
	}
}

func (h *NotificationHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.LogError(err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userId").(string)
	client := domain_notification.NewClient(userId, conn, h.hub, h.logger)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
