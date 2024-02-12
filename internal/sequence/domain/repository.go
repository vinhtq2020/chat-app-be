package domain

import "github.com/gin-gonic/gin"

type SequenceRepository interface {
	Next(c *gin.Context, module string) (int64, error)
	GetSequence(c *gin.Context, module string) (int64, error)
}
