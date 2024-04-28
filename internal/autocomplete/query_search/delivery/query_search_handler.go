package delivery

import (
	"go-service/internal/autocomplete/query_search/domain"
	"go-service/pkg/logger"
	"go-service/pkg/response"
	"net/http"
)

type QuerySearchHandler struct {
	service domain.QuerySearchService
	logger  *logger.Logger
}

func NewQuerySearchHandler(logger *logger.Logger, service domain.QuerySearchService) *QuerySearchHandler {
	return &QuerySearchHandler{
		logger:  logger,
		service: service,
	}
}

func (handler *QuerySearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) == 0 {
		response.Response(w, http.StatusBadRequest, nil)
		return
	}
	fileLogs := "log_query"
	handler.logger.LogInfo(q, &fileLogs)
	res, err := handler.service.Search(r.Context(), q)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response.Response(w, http.StatusOK, res)
}
