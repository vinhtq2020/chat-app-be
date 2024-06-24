package http

import (
	"context"
	"go-service/internal/friend/domain"
	"go-service/pkg/handler_fnc"
	"go-service/pkg/logger"
	"net/http"
)

func NewRoute(ctx context.Context, r *http.ServeMux, handler domain.FriendTransport, logger *logger.Logger) {
	handler_fnc.HandleWithSecurity(ctx, r, "", http.MethodPost, "/{userId}/add-friend/{friendId}", logger, true, handler.Create)
	handler_fnc.HandleWithSecurity(ctx, r, "", http.MethodPost, "/{userId}/{action}/{requestId}", logger, true, handler.Patch)
}
