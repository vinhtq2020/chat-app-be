package domain

type Hub interface {
	UnRegister(*Client)
	Register(*Client)
	Run()
	SendMessage(WsMsg)
}
