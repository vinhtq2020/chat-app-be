package domain

import "github.com/gin-gonic/gin"

type UserTransport interface {
	Load(e *gin.Context)
	Create(e *gin.Context)
	Patch(e *gin.Context)
	Delete(e *gin.Context)
}
