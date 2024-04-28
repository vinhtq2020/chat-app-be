package domain

import (
	"context"
	"go-service/internal/autocomplete/model"
)

type QuerySearchRepository interface {
	LoadTrie(ctx context.Context) (*model.Trie, error)
}
