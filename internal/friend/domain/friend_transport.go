package domain

import (
	"net/http"
)

type FriendTransport interface {
	Create(w http.ResponseWriter, r *http.Request)
	UpdateRequestStatus(w http.ResponseWriter, r *http.Request)
}
