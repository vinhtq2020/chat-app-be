package http

import (
	"go-service/internal/auth/domain"

	"github.com/gin-gonic/gin"
)

func NewRoute(r *gin.RouterGroup, handler domain.AuthTransport) {
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
}
