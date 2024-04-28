package domain

import "context"

type SequenceRepository interface {
	Next(ctx context.Context, module string) (int64, error)
	GetSequence(ctx context.Context, module string) (int64, error)
}
