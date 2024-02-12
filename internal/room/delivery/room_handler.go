package delivery

import (
	"go-service/internal/room/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type RoomHandler struct {
	service  domain.RoomService
	upgrader websocket.Upgrader
}

func NewRoomHandler(service domain.RoomService, upgrader websocket.Upgrader) *RoomHandler {
	return &RoomHandler{service: service, upgrader: upgrader}
}
func (r *RoomHandler) All(c *gin.Context) {
	room, err := r.service.All(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, room)
}

func reader(c *gin.Context, conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func (r *RoomHandler) ReadAndWriteMessage(c *gin.Context) {
	r.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := r.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	reader(c, ws)
}

func (r *RoomHandler) Load(c *gin.Context) {
	var room *domain.Room
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	room, err := r.service.Load(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if room == nil {
		c.JSON(http.StatusNotFound, 0)
		return
	}
	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) Create(c *gin.Context) {
	var room domain.Room
	err := c.Bind(&room)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Create(c, room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else if res == 0 {
		c.JSON(http.StatusNotFound, 0)
		return
	} else if res < 0 {
		c.JSON(http.StatusConflict, -1)
		return
	} else {
		c.JSON(http.StatusCreated, res)
	}
}

func (h *RoomHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if len(id) > 0 {
		res, err := h.service.Delete(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		} else if res == 0 {
			c.JSON(http.StatusNotFound, 0)
			return
		} else if res < 0 {
			c.JSON(http.StatusConflict, -1)
			return
		} else {
			c.JSON(http.StatusOK, res)
		}
	} else {
		c.JSON(http.StatusBadRequest, "id not found")
	}
}

func (h *RoomHandler) Patch(c *gin.Context) {
	var room domain.Room

	id := c.Param("id")
	if len(id) > 0 {
		room.Id = id
		err := c.Bind(&room)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.service.Patch(c, room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		} else if res == 0 {
			c.JSON(http.StatusNotFound, 0)
			return
		} else if res < 0 {
			c.JSON(http.StatusConflict, -1)
			return
		} else {
			c.JSON(http.StatusOK, res)
		}
	}

}
