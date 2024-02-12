package domain

import (
	"github.com/gin-gonic/gin"
)

type AuthRepository interface {
	Register(*gin.Context, UserLoginData) (int64, error)
	GetUserLoginData(e *gin.Context, email string) (*UserLoginData, error)
	InTransaction(e *gin.Context, tx func() (int64, error)) (int64, error)
}
