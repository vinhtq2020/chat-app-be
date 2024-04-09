package http

import (
	"go-service/internal/query_search/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRoute(r *gin.RouterGroup, handler domain.QuerySearchTransport) {
	r.Handle(http.MethodGet, "/search", handler.Search)
}
