package domain

import (
	"net/http"
)

type RoomTransport interface {
	All(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	HandleConnections(w http.ResponseWriter, r *http.Request)
	HandleMessages()
}
