package usecase

import (
	"context"
	"go-service/internal/autocomplete/query_search/domain"
)

type QuerySearchUsecase struct {
	repository domain.QuerySearchRepository
}

func NewQuerySearchUsecase(repository domain.QuerySearchRepository) *QuerySearchUsecase {
	return &QuerySearchUsecase{
		repository: repository,
	}
}

func (u *QuerySearchUsecase) Search(ctx context.Context, query string) ([]string, error) {
	trie, err := u.repository.LoadTrie(ctx)
	if err != nil {
		return nil, err
	}

	res := trie.Search(query)
	return res, nil
}
