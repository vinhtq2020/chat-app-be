package domain

import (
	"net/http"
)

type FriendTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
}
