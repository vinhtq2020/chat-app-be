package hub

import "go-service/internal/notification/domain"

type NotificationHub struct {
	clients    map[string]map[*domain.Client]bool
	broadcast  chan domain.WsMsg
	register   chan *domain.Client
	unregister chan *domain.Client
}

func NewNotificationHub() *NotificationHub {
	return &NotificationHub{
		clients:    map[string]map[*domain.Client]bool{},
		broadcast:  make(chan domain.WsMsg),
		register:   make(chan *domain.Client),
		unregister: make(chan *domain.Client),
	}
}

func (h *NotificationHub) Register(client *domain.Client) {
	h.register <- client
}

func (h *NotificationHub) UnRegister(client *domain.Client) {
	h.unregister <- client
}

func (h *NotificationHub) SendMessage(msg domain.WsMsg) {
	h.broadcast <- msg
}

func (h *NotificationHub) Run() {
	for {
		select {
		case register := <-h.register:
			h.clients[register.UserId] = map[*domain.Client]bool{
				register: true,
			}
		case unregister := <-h.unregister:

			if _, ok := h.clients[unregister.UserId]; ok {
				delete(h.clients[unregister.UserId], unregister)
				if len(h.clients[unregister.UserId]) == 0 {
					delete(h.clients, unregister.UserId)
				}
				close(unregister.Send)
			}
		case WsMsg := <-h.broadcast:
			for _, subscriber := range WsMsg.Data.(domain.Notification).Subscribers {
				if _, ok := h.clients[subscriber.Id]; ok {
					for client := range h.clients[subscriber.Id] {
						select {
						case client.Send <- WsMsg:
						default:
							close(client.Send)
							delete(h.clients[subscriber.Id], client)
							if len(h.clients[subscriber.Id]) == 0 {
								delete(h.clients, subscriber.Id)
							}
						}
					}
				}
			}
		}
	}
}
