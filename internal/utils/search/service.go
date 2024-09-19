package search

import "context"

type SearchService[F any] interface {
	Search(ctx context.Context, filter F) (list interface{}, total int64, err error)
}
