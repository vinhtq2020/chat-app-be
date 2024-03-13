package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

const (
	requiredErrorMsg    = "%s is required"
	minRequiredErrorMsg = "%s is not less than %s"
	maxRequiredErrorMsg = "%s is not greater than %s"
)

type ErrorMsg struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
}

type Validate interface {
	Validate(s interface{}) []ErrorMsg
}

type validate struct {
	validator *validator.Validate
}

func NewValidate(validator *validator.Validate) *validate {
	return &validate{validator: validator}
}

func (v *validate) Validate(s interface{}) (errMsgs []ErrorMsg) {
	err := v.validator.Struct(s)
	if err != nil {
		valErrs := err.(validator.ValidationErrors)
		for _, item := range valErrs {
			message := ""
			switch item.Tag() {
			case "required":
				message = fmt.Sprintf(requiredErrorMsg, item.Field())
			case "min":
				message = fmt.Sprintf(minRequiredErrorMsg, item.Field(), item.Param())
			case "max":
				message = fmt.Sprintf(maxRequiredErrorMsg, item.Field(), item.Param())
			case "required_without":
				message = fmt.Sprintf(requiredErrorMsg, item.Field())
			default:
				message = item.Error()
			}
			errMsgs = append(errMsgs, ErrorMsg{
				Code:    item.Tag(),
				Message: message,
				Field:   item.Field(),
			})
		}
	}
	return errMsgs
}
