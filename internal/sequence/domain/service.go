package domain

import "github.com/gin-gonic/gin"

type SequenceService interface {
	Next(c *gin.Context, module string) (int64, error)
}
