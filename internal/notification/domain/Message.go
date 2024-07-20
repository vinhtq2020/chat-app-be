package domain

import "github.com/gorilla/websocket"

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Client struct {
	Conn   *websocket.Conn
	UserId string
}
