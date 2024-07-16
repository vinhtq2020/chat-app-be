package notification_domain

import "net/http"

type NotificacationTransport interface {
	HandleConnections(w http.ResponseWriter, r *http.Request)
	HandleMessage()
}
