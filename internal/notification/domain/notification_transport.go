package domain

import "net/http"

type NotificacationTransport interface {
	ServeWs(w http.ResponseWriter, r *http.Request)
}
