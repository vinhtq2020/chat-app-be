package domain

import (
	"go-service/pkg/jwt"
	"go-service/pkg/validate"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(e *gin.Context, user UserLoginData) ([]validate.ErrorMsg, int64, error)
	LoginWithGoogle(e *gin.Context, email string) ([]validate.ErrorMsg, int64, *jwt.TokenData, error)
	Login(e *gin.Context, email string, password string, browser string, ipAdress string, deviceId string) ([]validate.ErrorMsg, *jwt.TokenData, error)
	Logout(e *gin.Context, userId string, browser string, ipAdress string, deviceId string) (int64, error)
}
