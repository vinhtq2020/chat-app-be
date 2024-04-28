package http

import (
	"go-service/internal/autocomplete/query_search/domain"
	"net/http"
)

func NewRoute(router *http.ServeMux, handler domain.QuerySearchTransport) {
	router.HandleFunc(http.MethodGet+" "+"/search", handler.Search)
}
