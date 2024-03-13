package delivery

import (
	"go-service/internal/auth/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
}

// Login implements domain.AuthTransport.
func (h *AuthHandler) Login(e *gin.Context) {
	var userLoginData domain.UserLoginData
	err := e.Bind(&userLoginData)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if userLoginData.Email != nil && userLoginData.Password != nil {
		listErrs, token, err := h.authService.Login(e, *userLoginData.Email, *userLoginData.Password)
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

// Register implements domain.AuthTransport.
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

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
