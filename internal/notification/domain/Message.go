package domain

import (
	"github.com/gorilla/websocket"
)

type NotificationMessage string

const (
	NOTIFIED NotificationMessage = "notified"
	UPDATED  NotificationMessage = "updated"
	READ     NotificationMessage = "read"
)

type Message struct {
	Name NotificationMessage `json:"name"`
	Data any                 `json:"data"`
}

type Client struct {
	Conn   *websocket.Conn
	UserId string
}
