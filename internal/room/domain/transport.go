package domain

import "github.com/gin-gonic/gin"

type RoomTransport interface {
	All(c *gin.Context)
	Load(c *gin.Context)
	Create(c *gin.Context)
	Patch(c *gin.Context)
	Delete(c *gin.Context)
}
