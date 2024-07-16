package http

import (
	"context"
	domain_notification "go-service/internal/notification/notification_domain"
	"go-service/pkg/handler_fnc"
	"go-service/pkg/logger"
	"net/http"
)

func Route(ctx context.Context, router *http.ServeMux, handler domain_notification.NotificacationTransport, logger *logger.Logger) {
	handler_fnc.HandleWithSecurity(ctx, router, "notification", http.MethodGet, "/ws", logger, true, handler.HandleConnections)
}
