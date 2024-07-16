package http

import (
	"go-service/internal/user/user_domain"
	"net/http"
)

func NewRoute(router *http.ServeMux, handler user_domain.UserTransport) {
	path := "/user"

	router.HandleFunc(path+"/search", handler.Search)
	router.HandleFunc(http.MethodPost+path+"/search", handler.Search)
}
