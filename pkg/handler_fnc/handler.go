package handler_fnc

import (
	"context"
	"fmt"
	"go-service/internal/configs"
	"go-service/pkg/jwt"
	"go-service/pkg/logger"
	"go-service/pkg/response"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"

	"github.com/gorilla/websocket"
)

func HandleWithSecurity(ctx context.Context, router *http.ServeMux, routerGroup string, httpMethod string, relativePath string, logger *logger.Logger, security bool, handlerFunc func(http.ResponseWriter, *http.Request)) {
	handlerFnc := func(w http.ResponseWriter, r *http.Request) {
		allow := ipFilter(ctx, r)
		if !allow {
			response.Response(w, http.StatusUnauthorized, "Unauthorized")
		} else if security {
			// error always return http.ErrNoCookie if not found cookie
			c, err := r.Cookie("accessToken")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userId, err := r.Cookie("userId")
			if err != nil {
				logger.LogError(err.Error(), nil)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			rctx := context.WithValue(r.Context(), "userId", userId.Value)
			r = r.WithContext(rctx)
			accessToken := c.Value
			// validate access token & refresh token & user id
			secretKey := ctx.Value("secretKey").(string)
			if len(secretKey) == 0 {
				logger.LogError("invalid secret key", nil)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			res, err := jwt.VerifyAccessToken(accessToken, secretKey)
			if err != nil {
				logger.LogError(err.Error(), nil)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			} else if res == 0 {
				http.Error(w, "Not Found", http.StatusNotFound)
			} else if res == -1 {
				response.Response(w, http.StatusInternalServerError, nil)
			} else if res == -2 {
				response.Response(w, http.StatusUnauthorized, "Unauthorized")
			} else if res == -3 {
				response.Response(w, http.StatusUnauthorized, "token is expired")

			}

			// // get refresh token is expired or not, if not, generate new access token
			// recoder := httptest.NewRecorder()
			// refreshToken(recoder, r)
			// // check refreshToken is done
			// if recoder.Code != 0 {
			// 	return
			// }

		}
		handlerFunc(w, r)

	}
	router.HandleFunc(httpMethod+" "+routerGroup+relativePath, handlerFnc)
}

func LogRequestHandler(h http.Handler, logger *logger.Logger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			logger.LogError(err.Error(), nil)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Println(fmt.Sprintf("%q", x))

		if websocket.IsWebSocketUpgrade(r) {
			h.ServeHTTP(w, r)
			return
		}

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		log.Println(fmt.Sprintf("%q", rec.Body))
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	}
	return http.HandlerFunc(fn)
}

func ipFilter(ctx context.Context, r *http.Request) bool {
	inbound := ctx.Value("inbound").(map[string]configs.Address)
	for _, v := range inbound {
		if v.Host == r.Host && v.Port == r.URL.Port() {
			return true
		}
	}
	return false
}
