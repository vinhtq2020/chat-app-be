package domain

import "github.com/gin-gonic/gin"

type AuthTransport interface {
	Register(e *gin.Context)
	Login(e *gin.Context)
	Logout(e *gin.Context)
}
