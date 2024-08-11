package search

import (
	"golang.org/x/net/context"
)

type searchUsecase[T any] struct {
	repo SearchRepository
}

func NewSearchService[T any](repo SearchRepository) *searchUsecase[T] {
	return &searchUsecase[T]{
		repo: repo,
	}
}

func (f *searchUsecase[T]) Search(ctx context.Context, filter SearchFilter) (interface{}, int64, error) {
	list := []T{}

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
