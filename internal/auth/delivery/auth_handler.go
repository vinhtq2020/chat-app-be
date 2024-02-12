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
func (*AuthHandler) Login(e *gin.Context) {
	panic("unimplemented")
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
		e.JSON(http.StatusCreated, res)
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
		e.JSON(http.StatusBadRequest, errs)

	} else if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
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
