package delivery

import (
	"encoding/json"
	"go-service/internal/room/domain"
	"go-service/pkg/logger"
	"go-service/pkg/response"
	"net/http"

	"github.com/gorilla/websocket"
)

const GROUP_LIMIT = 100

type RoomHandler struct {
	service   domain.RoomService
	upgrader  websocket.Upgrader
	clients   map[*websocket.Conn]bool
	limit     int8
	broadcast chan domain.Message
	logger    *logger.Logger
}

func NewRoomHandler(service domain.RoomService, upgrader websocket.Upgrader, logger *logger.Logger) *RoomHandler {
	broadcast := make(chan domain.Message)
	return &RoomHandler{service: service, upgrader: upgrader, limit: GROUP_LIMIT, broadcast: broadcast, logger: logger}
}
func (h *RoomHandler) All(w http.ResponseWriter, r *http.Request) {
	room, err := h.service.All(r.Context())
	if err != nil {
		h.logger.LogError(err.Error(), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Response(w, http.StatusOK, room)
}

func (h *RoomHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// handle connection
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.LogError(err.Error(), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer ws.Close()

	h.clients[ws] = true

	for {
		var msg domain.Message
		err = ws.ReadJSON(&msg)
		if err != nil {
			h.logger.LogError(err.Error(), nil)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			delete(h.clients, ws)
			break
		}

		h.broadcast <- msg
	}

}

func (h *RoomHandler) HandleMessages() {
	for {
		msg := <-h.broadcast

		for client := range h.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				h.logger.LogError(err.Error(), nil)
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}

func (h *RoomHandler) Load(w http.ResponseWriter, r *http.Request) {
	var room *domain.Room
	id := r.PathValue("id")
	if len(id) == 0 {
		http.Error(w, "missing required field", http.StatusBadRequest)
		return
	}
	room, err := h.service.Load(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if room == nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	response.Response(w, http.StatusOK, room)
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var room domain.Room

	err := json.NewDecoder(r.Body).Decode(&room)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.service.Create(r.Context(), room)
	if err != nil {
		response.Response(w, http.StatusInternalServerError, err.Error())
		return
	} else if res == 0 {
		response.Response(w, http.StatusNotFound, 0)
		return
	} else if res < 0 {
		response.Response(w, http.StatusConflict, -1)
		return
	} else {
		response.Response(w, http.StatusCreated, res)
	}
}

func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) > 0 {
		res, err := h.service.Delete(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if res == 0 {
			response.Response(w, http.StatusNotFound, 0)
			return
		} else if res < 0 {
			response.Response(w, http.StatusConflict, -1)
			return
		} else {
			response.Response(w, http.StatusOK, res)
		}
	} else {
		http.Error(w, "id not found", http.StatusBadRequest)
	}
}

func (h *RoomHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var room domain.Room

	id := r.PathValue("id")
	if len(id) > 0 {
		room.Id = id
		err := json.NewDecoder(r.Body).Decode(&room)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := h.service.Patch(r.Context(), room)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if res == 0 {
			response.Response(w, http.StatusNotFound, 0)
			return
		} else if res < 0 {
			response.Response(w, http.StatusConflict, -1)
			return
		} else {
			response.Response(w, http.StatusOK, res)
		}
	}

}
