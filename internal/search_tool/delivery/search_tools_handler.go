package delivery

import (
	"go-service/internal/search_tool/domain"
	"go-service/internal/search_tool/service"
	"go-service/internal/utils/search"
	"go-service/pkg/response"
	"net/http"
)

type SearchToolsHandler struct {
	service service.SearchToolService
}

func NewSearchToolsHandler(service service.SearchToolService) *SearchToolsHandler {
	return &SearchToolsHandler{
		service: service,
	}
}

func (h *SearchToolsHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := domain.SearchFilter{}
	err := search.Bind(r, &filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		response.Response(w, http.StatusBadRequest, nil)
		return
	}
	res, total, list, err := h.service.Search(r.Context(), userId, filter)
	if err != nil {
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
	} else if res > 0 {
		response.Response(w, http.StatusOK, domain.SearchResponse{
			List:  list,
			Total: total,
		})
	} else if res == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else if res == -1 {
		response.Response(w, http.StatusConflict, nil)
	} else if res == -2 {
		response.Response(w, http.StatusUnauthorized, nil)
	}
}
