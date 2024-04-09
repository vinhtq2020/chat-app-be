package route

import (
	"go-service/app"
	http_auth "go-service/internal/auth/delivery/http"
	http_query_search "go-service/internal/query_search/delivery/http"
	http_room "go-service/internal/room/delivery/http"
	http_user "go-service/internal/user/delivery/http"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine, app *app.App) {
	http_user.NewRoute(r.Group("/user"), app.User)
	http_room.NewRoute(r.Group("/room"), app.Room)
	http_auth.NewRoute(r.Group("/auth"), app.Auth)
	http_query_search.NewRoute(r.Group(""), app.QuerySearch)
}
