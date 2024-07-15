package http

import (
	"context"
	"go-service/internal/search_tool"
	"go-service/pkg/handler_fnc"
	"go-service/pkg/logger"
	"net/http"
)

func NewRoute(ctx context.Context, router *http.ServeMux, handler search_tool.SearchToolTransport, logger *logger.Logger) {
	path := ""
	handler_fnc.HandleWithSecurity(ctx, router, path, http.MethodGet, "/search", logger, true, handler.Search)
	handler_fnc.HandleWithSecurity(ctx, router, path, http.MethodPost, "/search", logger, true, handler.Search)

}
