package delivery

import (
	domain_notification "go-service/internal/notification/notification_domain"
	"go-service/pkg/logger"
	"net/http"

	"github.com/gorilla/websocket"
)

type NotificationHandler struct {
	clients   map[string]*websocket.Conn
	broadcast chan domain_notification.Notification
	upgrader  *websocket.Upgrader
	logger    *logger.Logger
}

func NewNotificationHandler(upgrader *websocket.Upgrader, logger *logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		clients:   map[string]*websocket.Conn{},
		broadcast: make(chan domain_notification.Notification),
		upgrader:  upgrader,
		logger:    logger,
	}
}

func (h *NotificationHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.LogError(err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer ws.Close()
	userId := r.Context().Value("userId").(string)
	h.clients[userId] = ws
	for {
		var noti domain_notification.Notification
		err := ws.ReadJSON(noti)
		if err != nil {
			h.logger.LogError(err.Error(), nil)
			delete(h.clients, userId)
			break
		}

		h.broadcast <- noti
	}
}

func (h *NotificationHandler) HandleMessage() {
	for {
		noti := <-h.broadcast
		for _, subscriber := range noti.Subscribers {
			conn, existed := h.clients[subscriber.Id]
			if !existed {
				continue
			}
			err := conn.WriteJSON(noti)
			if err != nil {
				h.logger.LogError(err.Error(), nil)
				conn.Close()
				delete(h.clients, subscriber.Id)
			}
		}
	}
}
