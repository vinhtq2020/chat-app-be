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
	handler_fnc.HandleWithSecurity(ctx, router, path, http.MethodPost, "/register", logger, false, handler.Register)
	handler_fnc.HandleWithSecurity(ctx, router, path, http.MethodPost, "/login", logger, false, handler.Login)
	handler_fnc.HandleWithSecurity(ctx, router, path, http.MethodGet, "/logout", logger, true, handler.Logout)
	handler_fnc.HandleWithSecurity(ctx, router, path, http.MethodGet, "/refresh", logger, false, handler.RefreshToken)
}
