package http

import (
	"go-service/internal/auth/domain"
	"go-service/pkg/handler_fnc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRoute(r *gin.RouterGroup, handler domain.AuthTransport) {
	r.Handle(http.MethodPost, "/register", handler.Register)
	handler_fnc.HandleWithSecurity(r, http.MethodPost, "/login", true, handler.Login)
	handler_fnc.HandleWithSecurity(r, http.MethodPost, "/logout", true, handler.Logout)
}
