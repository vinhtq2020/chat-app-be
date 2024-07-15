package domain

import "context"

type SearchToolRepository interface {
	Search(ctx context.Context, id string, filter SearchFilter) ([]SearchItem, error)
	Total(ctx context.Context, filter SearchFilter) (int64, error)
}
