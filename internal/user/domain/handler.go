package domain

import "github.com/gin-gonic/gin"

type UserTransport interface {
	Search(e *gin.Context)
	Load(e *gin.Context)
	Create(e *gin.Context)
	Patch(e *gin.Context)
	Delete(e *gin.Context)
}
