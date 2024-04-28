package domain

import "context"

type QuerySearchService interface {
	Search(ctx context.Context, query string) ([]string, error)
}
