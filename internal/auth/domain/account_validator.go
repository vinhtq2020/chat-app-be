package domain

import (
	"context"
	"go-service/pkg/validate"
)

type AccountValidator interface {
	ValidateLogin(ctx context.Context, email string) ([]validate.ErrorMsg, error)
	ValidateEmailGoogle(ctx context.Context, email string) ([]validate.ErrorMsg, error)
	ValidateRegister(ctx context.Context, user Account) ([]validate.ErrorMsg, error)
}
