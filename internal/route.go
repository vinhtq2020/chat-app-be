package route

import (
	"context"
	"go-service/app"
	http_auth "go-service/internal/auth/delivery/http"
	http_query_search "go-service/internal/autocomplete/query_search/delivery/http"
	http_room "go-service/internal/room/delivery/http"
	http_user "go-service/internal/user/delivery/http"
	"go-service/pkg/logger"
	"net/http"
)

func Route(ctx context.Context, router *http.ServeMux, app *app.App, logger *logger.Logger, refreshToken http.HandlerFunc) {
	http_user.NewRoute(router, app.User)
	http_room.NewRoute(router, app.Room)
	http_auth.NewRoute(ctx, router, app.Auth, logger)
	http_query_search.NewRoute(router, app.QuerySearch)
}
