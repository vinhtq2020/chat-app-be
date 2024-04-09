package delivery

import (
	"go-service/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QuerySearchHandler struct {
	logger *logger.Logger
}

func NewQuerySearchHandler(logger *logger.Logger) *QuerySearchHandler {
	return &QuerySearchHandler{
		logger: logger,
	}
}

func (handler *QuerySearchHandler) Search(e *gin.Context) {
	q := e.Query("q")
	if len(q) == 0 {
		e.JSON(http.StatusBadRequest, nil)
		return
	}
	fileLogs := "log_query"
	handler.logger.LogInfo(q, &fileLogs)
}
