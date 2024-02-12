package domain

import (
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	Create(e *gin.Context, user User) (int64, error)
}
