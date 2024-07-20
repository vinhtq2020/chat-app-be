package domain

type WsMsg struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}
