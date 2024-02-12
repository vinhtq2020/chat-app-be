package domain

import "github.com/gin-gonic/gin"

type FriendTransport interface {
	Create(e *gin.Context)
	Update(e *gin.Context)
	Patch(e *gin.Context)
	Delete(e *gin.Context)
}
