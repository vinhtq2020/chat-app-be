package search_tool

import (
	"go-service/internal/search_tool/delivery"
	"go-service/internal/search_tool/repository"
	"go-service/internal/search_tool/service"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"
	"net/http"

	"gorm.io/gorm"
)

type SearchToolTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
}

func NewSearchToolTransport(db *gorm.DB, logger *logger.Logger, toArray pq.Array) SearchToolTransport {
	repository := repository.NewSearchToolsRepository(db, toArray, logger)
	sv := service.NewSearchToolService(repository)
	handler := delivery.NewSearchToolsHandler(sv)
	return handler
}
