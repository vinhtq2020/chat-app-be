package service

import (
	"go-service/internal/search/domain"

	"golang.org/x/net/context"
)

type SearchService[T any] struct {
	repo domain.SearchRepository
}

func NewSearchService[T any](repo domain.SearchRepository) *SearchService[T] {
	return &SearchService[T]{
		repo: repo,
	}
}

func (f *SearchService[T]) Search(ctx context.Context, filter domain.SearchFilter) (interface{}, int64, error) {
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
