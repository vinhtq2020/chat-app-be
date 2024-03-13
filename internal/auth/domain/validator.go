package domain

import (
	"go-service/pkg/validate"

	"github.com/gin-gonic/gin"
)

type UserLoginDataValidator interface {
	ValidateLogin(e *gin.Context, email string) ([]validate.ErrorMsg, error)
	ValidateEmailGoogle(e *gin.Context, email string) ([]validate.ErrorMsg, error)
	ValidateRegister(e *gin.Context, user UserLoginData) ([]validate.ErrorMsg, error)
}
