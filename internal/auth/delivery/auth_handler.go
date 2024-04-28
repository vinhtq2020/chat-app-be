package delivery

import (
	"encoding/json"
	"go-service/internal/auth/domain"
	"go-service/pkg/jwt"
	"go-service/pkg/logger"
	"go-service/pkg/response"
	useragent "go-service/pkg/user_agent"
	"net/http"
	"time"
)

type AuthHandler struct {
	authService domain.AuthService
	logger      *logger.Logger
}

func NewAuthHandler(authService domain.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var userLoginData domain.Account
	err := json.NewDecoder(r.Body).Decode(&userLoginData)
	userAgent := r.UserAgent()
	deviceId := r.Header.Get("Device-ID")
	ipAddress := r.Header.Get("X-Forwarded-For")

	browser, _ := useragent.GetDeviceInfo(userAgent)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(deviceId) == 0 || len(browser) == 0 || len(ipAddress) == 0 {
		http.Error(w, "missing required Header", http.StatusBadRequest)
		return
	}

	if userLoginData.Email != nil && userLoginData.Password != nil {
		listErrs, token, err := h.authService.Login(r.Context(), *userLoginData.Email, *userLoginData.Password, browser, ipAddress, deviceId)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else if len(listErrs) > 0 {
			response.Response(w, http.StatusUnprocessableEntity, listErrs)
		} else if token == nil {
			http.Error(w, "email is not existed", http.StatusNotFound)
		} else {
			expiresAt := time.Now().Add(jwt.AccessTokenDuration)
			http.SetCookie(w, &http.Cookie{
				Name:     "accessToken",
				Value:    token.AccessToken,
				Expires:  expiresAt,
				Secure:   true,
				HttpOnly: true,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "userId",
				Value:    token.UserId,
				Secure:   true,
				HttpOnly: true,
			})
			response.Response(w, http.StatusOK, 1)
		}
		return
	}

	http.Error(w, "email or password cannot be null", http.StatusBadRequest)

}

func (h *AuthHandler) LoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	var userLoginData domain.Account
	err := json.NewDecoder(r.Body).Decode(&userLoginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listErrs, res, err := h.authService.Register(r.Context(), userLoginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if len(listErrs) > 0 {
		response.Response(w, http.StatusUnprocessableEntity, listErrs)
	} else {
		response.Response(w, http.StatusOK, res)

	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userLoginData domain.Account
	err := json.NewDecoder(r.Body).Decode(&userLoginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	errs, res, err := h.authService.Register(r.Context(), userLoginData)
	if len(errs) > 0 {
		response.Response(w, http.StatusUnprocessableEntity, errs)
	} else if err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
	} else if res > 0 {
		response.Response(w, http.StatusOK, res)
	} else if res == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if res == -1 {
		response.Response(w, http.StatusInternalServerError, nil)
	} else if res == -2 {
		response.Response(w, http.StatusUnauthorized, nil)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userId := r.Header["User-Id"][0]
	userAgent := r.UserAgent()
	deviceId := r.Header[http.CanonicalHeaderKey("Device-ID")][0]
	ipAddress := r.Header.Get("X-Forwarded-For")
	browser, _ := useragent.GetDeviceInfo(userAgent)

	if len(browser) == 0 || len(deviceId) == 0 || len(ipAddress) == 0 || len(userId) == 0 {
		http.Error(w, "Missing required headers", http.StatusBadRequest)
		return
	}

	res, err := h.authService.Logout(r.Context(), userId, browser, ipAddress, deviceId)
	if err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
	} else if res > 0 {
		response.Response(w, http.StatusOK, res)
	} else if res == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if res == -1 {
		response.Response(w, http.StatusInternalServerError, nil)
	} else if res == -2 {
		response.Response(w, http.StatusUnauthorized, nil)
	}
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	deviceId := r.Header.Get("Device-ID")
	ipAddress := r.Header.Get("X-Forwarded-For")
	userId, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	res, newToken, err := h.authService.RefreshToken(r.Context(), userId.Value, userAgent, ipAddress, deviceId)
	if err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
	} else if res > 0 {
		response.Response(w, http.StatusOK, res)
	} else if res == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if res == -1 {
		response.Response(w, http.StatusInternalServerError, nil)
	} else if res == -2 {
		response.Response(w, http.StatusUnauthorized, nil)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    newToken,
		Expires:  time.Now().Add(jwt.AccessTokenDuration),
		HttpOnly: true,
		Secure:   true,
	})
	// response.Response(w, http.StatusCreated)
}
