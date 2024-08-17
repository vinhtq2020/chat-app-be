package repository

import (
	"context"
	"fmt"
	"go-service/internal/search_tool/domain"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

// Search use for page search not autocomplete guild searcj
type SearchToolRepository struct {
	DB         *gorm.DB
	logger     *logger.Logger
	toArray    pq.Array
	buildParam func(int) string
}

func NewSearchToolsRepository(DB *gorm.DB, buildParam func(int) string, toArray pq.Array, logger *logger.Logger) domain.SearchToolRepository {
	return &SearchToolRepository{
		toArray:    toArray,
		DB:         DB,
		logger:     logger,
		buildParam: buildParam,
	}
}

func (r *SearchToolRepository) Search(ctx context.Context, id string, filter domain.SearchFilter) ([]domain.SearchItem, error) {
	var res []domain.SearchItem
	qr := `select a.id, a.user_name, a.avatar_url,
					case 
						when a.id = %s then NULL
						when b.user_id1 is not NULL and b.user_id2 is not NULL then b.status
						else 'none'
					end as friend_status
					from users a left join friends b on (a.id = b.user_id1 or a.id = b.user_id2)
					where a.user_name like CONCAT('%%',%s::text,'%%')`

	stmt := fmt.Sprintf(qr, r.buildParam(1), r.buildParam(2))
	err := sql.QueryWithArray(r.DB, &res, stmt, r.toArray, r.logger, id, *filter.Q)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
	}
	return res, err
}

// Total implements domain.SearchToolRepository.
func (r *SearchToolRepository) Total(ctx context.Context, filter domain.SearchFilter) (int64, error) {
	panic("unimplemented")
}
