package validator

import (
	"go-service/internal/auth/domain"
	"go-service/pkg/validate"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthValidator struct {
	validate *validator.Validate
}

func NewAuthValidator(validate *validator.Validate) *AuthValidator {
	return &AuthValidator{validate: validate}
}

func (v *AuthValidator) ValidateLogin(*gin.Context, domain.UserLoginData) ([]validate.ErrorMsg, error) {
	panic("")
}

func (v *AuthValidator) ValidateEmailGoogle(e *gin.Context, email string) ([]validate.ErrorMsg, error) {
	errMsgs := []validate.ErrorMsg{}
	_, err := mail.ParseAddress(email)
	if err != nil {
		errMsgs = append(errMsgs, validate.ErrorMsg{
			Message: err.Error(),
		})
	}
	emailParts := strings.Split(email, "@")
	if len(emailParts) <= 0 || emailParts[1] != "gmail.com" {
		errMsgs = append(errMsgs, validate.ErrorMsg{Message: "email is not valid"})
	}
	return nil, nil
}
