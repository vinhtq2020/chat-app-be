package http

import (
	"context"
	domain_notification "go-service/internal/notification/domain"
	"go-service/pkg/handler_fnc"
	"go-service/pkg/logger"
	"net/http"
)

func NewRoute(ctx context.Context, router *http.ServeMux, handler domain_notification.NotificacationTransport, logger *logger.Logger) {
	handler_fnc.HandleWithSecurity(ctx, router, "/notification", http.MethodGet, "/ws", logger, true, handler.ServeWs)
}
