package domain

import (
	"context"
)

type SearchRepository interface {
	Search(ctx context.Context, result interface{}, filter SearchFilter) error
	Total(ctx context.Context) (int64, error)
}
