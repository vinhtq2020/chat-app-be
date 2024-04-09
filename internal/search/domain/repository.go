package domain

import "github.com/gin-gonic/gin"

type SearchRepository interface {
	Search(e *gin.Context, result interface{}, filter SearchFilter) error
	Total(e *gin.Context) (int64, error)
}
