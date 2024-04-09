package delivery

import (
	domain_search "go-service/internal/search/domain"
	"go-service/internal/user/domain"
	"go-service/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (u *UserHandler) Search(e *gin.Context) {
	var filter domain_search.SearchFilter
	err := e.Bind(&filter)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return
	}
	list, total, err := u.searchService.Search(e, filter)
	if err != nil {
		e.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	e.JSON(http.StatusOK, model.SearchResult{
		List: list, Total: total,
	})
}

func (*UserHandler) Create(e *gin.Context) {
	var user domain.User

	err := e.Bind(&user)

	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return
	}

}

func (*UserHandler) Load(e *gin.Context) {
	panic("unimplemented")
}

func (*UserHandler) Delete(e *gin.Context) {
	panic("unimplemented")
}

// Patch implements domain.UserTransport.
func (*UserHandler) Patch(e *gin.Context) {
	panic("unimplemented")
}
