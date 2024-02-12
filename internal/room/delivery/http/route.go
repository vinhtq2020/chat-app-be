package http

import (
	"go-service/internal/room/domain"

	"github.com/gin-gonic/gin"
)

func NewRoute(r *gin.RouterGroup, handler domain.RoomTransport) {
	r.GET("", handler.All)
	r.GET("/:id", handler.Load)
	r.POST("", handler.Create)
	r.PATCH("/:id", handler.Patch)
	r.DELETE("/:id", handler.Delete)
}
