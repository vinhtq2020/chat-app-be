package http

import (
	user_domain "go-service/internal/user/domain"
	"net/http"
)

func NewRoute(router *http.ServeMux, handler user_domain.UserTransport) {
	path := "/user"

	router.HandleFunc(path+"/search", handler.Search)
	router.HandleFunc(http.MethodPost+path+"/search", handler.Search)
}
