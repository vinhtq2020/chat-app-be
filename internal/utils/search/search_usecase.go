package search

import (
	"golang.org/x/net/context"
)

type searchUsecase[F any] struct {
	repo SearchRepository[F]
}

func NewSearchService[F any](repo SearchRepository[F]) *searchUsecase[F] {
	return &searchUsecase[F]{
		repo: repo,
	}
}

func (f *searchUsecase[F]) Search(ctx context.Context, filter F) (interface{}, int64, error) {
	list := []interface{}{}

	err := f.repo.Search(ctx, &list, filter)
	if err != nil {
		return nil, 0, err
	}
	total, err := f.repo.Total(ctx)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
