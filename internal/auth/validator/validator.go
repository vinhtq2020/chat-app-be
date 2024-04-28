package validator

import (
	"context"
	"fmt"
	"go-service/internal/auth/domain"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/validate"
	"net/mail"
	"strings"

	"gorm.io/gorm"
)

type AuthValidator struct {
	table    string
	db       *gorm.DB
	validate validate.Validate
	toArray  pq.Array
}

func NewAuthValidator(db *gorm.DB, table string, validate validate.Validate, toArray pq.Array) *AuthValidator {
	return &AuthValidator{validate: validate, db: db, toArray: toArray, table: table}
}

func (v *AuthValidator) ValidateLogin(ctx context.Context, email string) ([]validate.ErrorMsg, error) {
	errMsgs := []validate.ErrorMsg{}
	// qr := fmt.Sprintf(`select email from %s u where u.user_name = $1`, v.table)
	// res := []domain.UserLoginData{}

	// // Login with phone
	// err := sql.QueryWithArray(v.db, &res, qr, v.toArray, email)
	// if err != nil {
	// 	return nil, err
	// }

	// if len(res) == 0 {
	// 	errMsgs = append(errMsgs, validate.ErrorMsg{
	// 		Message: "this account is not existed. Please try another account",
	// 	})
	// }
	return errMsgs, nil
}

func (v *AuthValidator) ValidateRegister(ctx context.Context, user domain.Account) ([]validate.ErrorMsg, error) {
	errMsgs := []validate.ErrorMsg{}
	errs := v.validate.Validate(user)
	if len(errs) > 0 {
		errMsgs = append(errMsgs, errs...)
	}
	qr := fmt.Sprintf(`
						SELECT 'username' AS existing_field, user_name AS existing_value
						FROM %s
						WHERE user_name = $2
						UNION
						SELECT 'email' AS existing_field, email AS existing_value
						FROM %s
						WHERE email = $1`, v.table, v.table)
	res := []domain.ExistingField{}
	err := sql.QueryWithArray(v.db, &res, qr, v.toArray, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		if v.ExistingField == "username" {
			errMsgs = append(errMsgs, validate.ErrorMsg{
				Message: "username already existed",
				Code:    "username",
				Field:   "username",
			})
		} else if v.ExistingField == "email" {
			errMsgs = append(errMsgs, validate.ErrorMsg{
				Message: "email already existed",
				Code:    "email",
				Field:   "email",
			})
		}
	}
	return errMsgs, nil
}

func (v *AuthValidator) ValidateEmailGoogle(ctx context.Context, email string) ([]validate.ErrorMsg, error) {
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
