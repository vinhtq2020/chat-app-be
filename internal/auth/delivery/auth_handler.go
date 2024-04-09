package delivery

import (
	"go-service/internal/auth/domain"
	useragent "go-service/pkg/user_agent"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(e *gin.Context) {
	var userLoginData domain.UserLoginData
	err := e.Bind(&userLoginData)
	userAgent := e.Request.UserAgent()
	deviceId := e.Request.Header.Get("DeviceID")
	ipAddress := e.ClientIP()
	browser, _ := useragent.GetDeviceInfo(userAgent)

	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if len(deviceId) == 0 || len(browser) == 0 || len(ipAddress) == 0 {
		e.JSON(http.StatusBadRequest, "missing required Header")
		return
	}

	if userLoginData.Email != nil && userLoginData.Password != nil {
		listErrs, token, err := h.authService.Login(e, *userLoginData.Email, *userLoginData.Password, browser, ipAddress, deviceId)
		if err != nil {
			e.JSON(http.StatusInternalServerError, "Internal Server Error")
		} else if len(listErrs) > 0 {
			e.JSON(http.StatusUnprocessableEntity, listErrs)
		} else if token == nil {
			e.JSON(http.StatusNotFound, "email is not existed")
		} else {
			e.JSON(http.StatusOK, token)
		}
	} else {
		e.JSON(http.StatusBadRequest, "email or password cannot be null")
	}
}

func (h *AuthHandler) LoginWithGoogle(e *gin.Context) {
	var userLoginData domain.UserLoginData
	err := e.Bind(&userLoginData)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return
	}
	listErrs, res, err := h.authService.Register(e, userLoginData)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
	} else if len(listErrs) > 0 {
		e.JSON(http.StatusBadRequest, listErrs)
	} else {
		e.JSON(http.StatusOK, res)
	}

}

func (h *AuthHandler) Register(e *gin.Context) {
	var userLoginData domain.UserLoginData
	err := e.Bind(&userLoginData)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return
	}
	errs, res, err := h.authService.Register(e, userLoginData)
	if len(errs) > 0 {
		e.JSON(http.StatusUnprocessableEntity, errs)

	} else if err != nil {
		e.JSON(http.StatusInternalServerError, "Internal Server Error")
	} else if res == 0 {
		e.JSON(http.StatusNotFound, res)
	} else if res < 0 {
		e.JSON(http.StatusConflict, res)
	} else {
		e.JSON(http.StatusCreated, res)
	}
}

func (h *AuthHandler) Logout(e *gin.Context) {
	userId := e.Request.Header.Get("UserId")
	userAgent := e.Request.UserAgent()
	deviceId := e.Request.Header.Get("deviceId")
	ipAddress := e.ClientIP()
	browser, _ := useragent.GetDeviceInfo(userAgent)

	if len(browser) == 0 || len(deviceId) == 0 || len(ipAddress) == 0 || len(userId) == 0 {
		e.JSON(http.StatusBadRequest, "Missing required headers")
		return
	}

	res, err := h.authService.Logout(e, userId, browser, ipAddress, deviceId)
	if res < 0 {
		e.JSON(http.StatusInternalServerError, "Internal Sever Error")
		return
	} else if res > 0 {
		e.JSON(http.StatusOK, res)
	} else {
		e.JSON(http.StatusNotFound, err.Error())
	}
}
