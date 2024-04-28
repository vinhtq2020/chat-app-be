package delivery

import (
	"encoding/json"
	domain_search "go-service/internal/search/domain"
	"go-service/internal/user/domain"
	"go-service/pkg/convert"
	"go-service/pkg/model"
	"net/http"
)

type UserHandler struct {
	service       domain.UserService
	searchService domain_search.SearchService
}

func NewUserHandler(service domain.UserService, searchService domain_search.SearchService) *UserHandler {
	return &UserHandler{
		service:       service,
		searchService: searchService,
	}
}

func (u *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	var filter domain_search.SearchFilter

	if r.Method == http.MethodGet {
		queryParams := r.URL.Query()
		filterMap := convert.ToMapOmitEmpty(filter)
		for k, v := range queryParams {
			if len(v) > 0 {
				filterMap[k] = v[0]
			}
		}
		jsonBody, err := json.Marshal(filterMap)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(jsonBody, &filter)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		err := json.NewDecoder(r.Body).Decode(&filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	list, total, err := u.searchService.Search(r.Context(), filter)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(model.SearchResult{
		List: list, Total: total,
	})

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonBody)

}

func (*UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (*UserHandler) Load(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (*UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// Patch implements domain.UserTransport.
func (*UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
