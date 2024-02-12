package http

import (
	"go-service/internal/user/domain"

	"github.com/gin-gonic/gin"
)

func NewRoute(r *gin.RouterGroup, handler domain.UserTransport) {
}
