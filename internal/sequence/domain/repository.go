package domain

import (
	"context"

	"gorm.io/gorm"
)

type SequenceRepository interface {
	Next(ctx context.Context, module string) (int64, error)
	GetSequence(ctx context.Context, module string) (int64, error)
	InTransaction(ctx context.Context, ex func(ctx context.Context, tx *gorm.DB) (int64, error)) (int64, error)
}
