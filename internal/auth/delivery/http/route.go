package http

import (
	"context"
	"go-service/internal/auth/domain"
	"go-service/pkg/handler_fnc"
	"go-service/pkg/logger"
	"net/http"
)

func NewRoute(ctx context.Context, router *http.ServeMux, handler domain.AuthTransport, logger *logger.Logger) {
	path := "/auth"
	handler_fnc.HandleWithSecurity(ctx, router, path, http.MethodPost, "/register", handler.RefreshToken, logger, false, handler.Register)
	handler_fnc.HandleWithSecurity(ctx, router, path, http.MethodPost, "/login", handler.RefreshToken, logger, false, handler.Login)
	handler_fnc.HandleWithSecurity(ctx, router, path, http.MethodPost, "/logout", handler.RefreshToken, logger, true, handler.Logout)
}
