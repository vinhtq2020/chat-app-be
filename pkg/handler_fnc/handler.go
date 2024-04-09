package handler_fnc

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleWithSecurity(r *gin.RouterGroup, httpMethod string, relativePath string, security bool, handlers ...gin.HandlerFunc) {
	handlerFncs := []gin.HandlerFunc{}
	for _, hf := range handlers {
		if security {
			handlerFnc := func(e *gin.Context) {
				accessToken := e.Request.Header.Get("Access-Token")
				refreshToken := e.Request.Header.Get("Refresh-Token")
				userId := e.Request.Header.Get("UserId")
				if len(refreshToken) > 0 && len(accessToken) > 0 && len(userId) > 0 {
					e.JSON(http.StatusUnauthorized, nil)
					return
				}
				hf(e)
			}

			handlerFncs = append(handlerFncs, handlerFnc)
		} else {
			handlerFncs = append(handlerFncs, hf)
		}
	}

	r.Handle(httpMethod, relativePath, handlerFncs...)
}
