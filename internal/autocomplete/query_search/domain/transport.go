package domain

import (
	"net/http"
)

type QuerySearchTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
}
