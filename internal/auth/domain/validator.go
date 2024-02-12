package domain

import (
	"go-service/pkg/validate"

	"github.com/gin-gonic/gin"
)

type UserLoginDataValidator interface {
	ValidateLogin(*gin.Context, UserLoginData) ([]validate.ErrorMsg, error)
	ValidateEmailGoogle(e *gin.Context, email string) ([]validate.ErrorMsg, error)
}
