package delivery

import (
	"go-service/internal/user/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Create implements domain.UserTransport.
func (*UserHandler) Create(e *gin.Context) {
	var user domain.User

	err := e.Bind(&user)

	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return
	}

}

// Load implements domain.UserTransport.
func (*UserHandler) Load(e *gin.Context) {
	panic("unimplemented")
}

// Delete implements domain.UserTransport.
func (*UserHandler) Delete(e *gin.Context) {
	panic("unimplemented")
}

// Patch implements domain.UserTransport.
func (*UserHandler) Patch(e *gin.Context) {
	panic("unimplemented")
}
