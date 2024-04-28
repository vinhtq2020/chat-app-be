package domain

import (
	"context"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Load(ctx context.Context, email string) (*Account, error)
	Exist(ctx context.Context, email string) (int64, error)
	Insert(ctx context.Context, account Account) (int64, error)
	InTransaction(ctx context.Context, ex func(tx *gorm.DB) (int64, error)) (int64, error)
}
