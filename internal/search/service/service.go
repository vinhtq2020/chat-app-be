package service

import (
	"go-service/internal/search/domain"

	"github.com/gin-gonic/gin"
)

type SearchService[T any] struct {
	repo domain.SearchRepository
}

func NewSearchService[T any](repo domain.SearchRepository) *SearchService[T] {
	return &SearchService[T]{
		repo: repo,
	}
}

func (f *SearchService[T]) Search(e *gin.Context, filter domain.SearchFilter) (interface{}, int64, error) {
	list := []T{}

	err := f.repo.Search(e, &list, filter)
	if err != nil {
		return nil, 0, err
	}
	total, err := f.repo.Total(e)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
