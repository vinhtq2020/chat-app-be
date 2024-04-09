package http

import (
	"go-service/internal/user/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRoute(r *gin.RouterGroup, handler domain.UserTransport) {
	r.Handle(http.MethodPost, "/search", handler.Search)
	r.Handle(http.MethodGet, "/search", handler.Search)
}
