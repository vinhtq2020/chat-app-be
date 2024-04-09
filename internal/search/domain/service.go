package domain

import "github.com/gin-gonic/gin"

type SearchService interface {
	Search(ctx *gin.Context, filter SearchFilter) (list interface{}, total int64, err error)
}
