package search

import (
	"context"
)

type SearchRepository[F any] interface {
	Search(ctx context.Context, result interface{}, filter F) error
	Total(ctx context.Context) (int64, error)
}
