package domain

import (
	"context"
)

type SequenceService interface {
	Next(ctx context.Context, module string) (int64, error)
}
