package service

import (
	"context"
	"go-service/internal/search_tool/domain"
)

type SearchToolService struct {
	repository domain.SearchToolRepository
}

func NewSearchToolService(repository domain.SearchToolRepository) SearchToolService {
	return SearchToolService{
		repository: repository,
	}
}

func (sv *SearchToolService) Search(ctx context.Context, filter domain.SearchFilter) (code int64, total int64, list []domain.SearchItem, err error) {
	userId := ctx.Value("userId")
	if userId == nil {
		userId = ""
	}
	res, err := sv.repository.Search(ctx, userId.(string), filter)
	if err != nil {
		return -1, 0, res, err
	}
	return 1, 0, res, nil
}
