package domain

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserLoginData(e *gin.Context, email string) (*UserLoginData, error)
	Exist(e *gin.Context, email string) (int64, error)
	AddRefreshToken(e *gin.Context, refreshToken RefreshToken) (int64, error)
	RemoveRefreshToken(e *gin.Context, userId string, ipAddress string, deviceId string, browser string) (int64, error)
	Register(*gin.Context, UserLoginData) (int64, error)
	InTransaction(e *gin.Context, ex func(tx *gorm.DB) (int64, error)) (int64, error)
}
