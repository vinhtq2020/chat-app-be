package delivery

import (
	"go-service/internal/notification/domain"
	"go-service/internal/utils/search"
	"go-service/pkg/logger"
	"go-service/pkg/response"
	"net/http"

	"github.com/gorilla/websocket"
)

type NotificationHandler struct {
	service    domain.NotificationService
	clients    map[string]map[domain.Client]bool
	broastcast chan domain.Message
	upgrader   *websocket.Upgrader
	logger     *logger.Logger
}

func NewNotificationHandler(upgrader *websocket.Upgrader, broastcast chan domain.Message, logger *logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		upgrader:   upgrader,
		logger:     logger,
		broastcast: broastcast,
		clients:    make(map[string]map[domain.Client]bool),
	}
}

func (h *NotificationHandler) Search(w http.ResponseWriter, r *http.Request) {
	var filter domain.NotificationFilter
	err := search.Bind(r, &filter)
	if err != nil {
		response.Response(w, http.StatusBadRequest, nil)
		return
	}

	filter.Visible = true
	h.service.Search(r.Context(), filter)
}

func (h *NotificationHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.LogError(err.Error(), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer func() {
		conn.Close()
	}()
	userId := r.Context().Value("userId").(string)

	if _, exist := h.clients[userId]; exist {
		h.clients[userId][domain.Client{
			Conn:   conn,
			UserId: userId,
		}] = true
	} else {
		h.clients[userId] = map[domain.Client]bool{
			{
				Conn:   conn,
				UserId: userId,
			}: true,
		}
	}

	for {
		var ms domain.Message
		err := conn.ReadJSON(&ms)
		if err != nil {
			h.logger.LogError(err.Error(), nil)
			delete(h.clients, userId)
			break
		}

		h.broastcast <- ms
	}
}

func (h *NotificationHandler) HandleMessages() {
	for {
		msg := <-h.broastcast
		switch msg.Name {
		case "notified":
			for _, subscriber := range msg.Data.(domain.Notification).Subscribers {
				if _, ok := h.clients[subscriber.Id]; ok {
					for client := range h.clients[subscriber.Id] {
						err := client.Conn.WriteJSON(msg)
						if err != nil {
							client.Conn.Close()
							delete(h.clients, client.UserId)
						}
					}
				}
			}
		default:
		}
	}
}
