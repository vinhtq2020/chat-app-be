package http

import (
	"go-service/internal/room/domain"
	"net/http"
)

func NewRoute(r *http.ServeMux, handler domain.RoomTransport) {
	path := "/room"
	r.HandleFunc(http.MethodGet+" "+path, handler.All)
	r.HandleFunc(http.MethodGet+" "+path+"/{id}", handler.Load)
	r.HandleFunc(http.MethodPost+" "+path, handler.Create)
	r.HandleFunc(http.MethodPatch+" "+path+"/{id}", handler.Patch)
	r.HandleFunc(http.MethodDelete+" "+path+"/{id}", handler.Delete)
}
