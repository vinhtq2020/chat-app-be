package http

import (
	"go-service/internal/user/domain"
	"net/http"
)

func NewRoute(router *http.ServeMux, handler domain.UserTransport) {
	path := "/user"

	router.HandleFunc(path+"/search", handler.Search)
	router.HandleFunc(http.MethodPost+path+"/search", handler.Search)
}
