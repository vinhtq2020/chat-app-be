package domain

import "context"

type SearchService interface {
	Search(ctx context.Context, filter SearchFilter) (list interface{}, total int64, err error)
}
