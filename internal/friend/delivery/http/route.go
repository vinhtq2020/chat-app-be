package http

import (
	"context"
	friend_domain "go-service/internal/friend/domain"
	"go-service/pkg/handler_fnc"
	"go-service/pkg/logger"
	"net/http"
)

func NewRoute(ctx context.Context, r *http.ServeMux, handler friend_domain.FriendTransport, logger *logger.Logger) {
	friend := "/friend"
	handler_fnc.HandleWithSecurity(ctx, r, friend, http.MethodPost, "/{friendId}", logger, true, handler.Create)
	handler_fnc.HandleWithSecurity(ctx, r, friend, http.MethodPatch, "/{friendId}/{action}", logger, true, handler.Patch)
}
