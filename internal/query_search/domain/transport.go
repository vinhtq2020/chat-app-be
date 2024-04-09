package domain

import "github.com/gin-gonic/gin"

type QuerySearchTransport interface {
	Search(e *gin.Context)
}
