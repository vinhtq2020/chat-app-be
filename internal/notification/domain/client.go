package domain

import (
	"go-service/pkg/logger"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	Send   chan WsMsg
	hub    Hub
	logger *logger.Logger
	UserId string
}

func NewClient(userId string, conn *websocket.Conn, hub Hub, logger *logger.Logger) *Client {
	return &Client{
		conn:   conn,
		Send:   make(chan WsMsg),
		hub:    hub,
		logger: logger,
		UserId: userId,
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {

		case notification, ok := <-c.Send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteJSON(notification)
			if err != nil {
				c.logger.LogError(err.Error(), nil)
			}

		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.UnRegister(c)
		c.conn.Close()
	}()

	for {
		var msg WsMsg
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.LogError(err.Error(), nil)

			}
			break
		}
		switch msg.Name {
		case "connect":
			c.hub.SendMessage(WsMsg{
				Name: "connect",
				Data: "is connected",
			})
		case "disconnect":
			c.hub.SendMessage(WsMsg{
				Name: "disconnect",
				Data: "is disconnected",
			})
			c.hub.UnRegister(c)

		}

	}
}
