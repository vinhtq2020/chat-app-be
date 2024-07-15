package repository

import (
	"context"
	"fmt"
	"go-service/internal/search_tool/domain"
	"go-service/internal/utils/search"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

// Search use for page search not autocomplete guild searcj
type SearchToolRepository struct {
	DB      *gorm.DB
	logger  *logger.Logger
	toArray pq.Array
}

func NewSearchToolsRepository(DB *gorm.DB, toArray pq.Array, logger *logger.Logger) domain.SearchToolRepository {
	return &SearchToolRepository{
		toArray: toArray,
		DB:      DB,
		logger:  logger,
	}
}

func (r *SearchToolRepository) Search(ctx context.Context, id string, filter domain.SearchFilter) ([]domain.SearchItem, error) {
	var res []domain.SearchItem
	selectClause := `select a.id, a.user_name, a.avatar_url,
					case 
						when b.user_id1 <> NULL and b.user_id2 <> NULL then b.status
						when a.id = %s then NULL
						else 'none'
					end as friend_status
					from users a left join friends b on (a.id = user_id1 or a.id = user_id2) 
					%s and a.user_name like CONCAT('%%',%s::text,'%%')`
	filterClause := search.BuildFilter(filter.SearchFilter)
	stmt := fmt.Sprintf(selectClause, sql.BuildParam(1), filterClause, sql.BuildParam(2))
	err := sql.QueryWithArray(r.DB, &res, stmt, r.toArray, id, filter.Q)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
	}
	return res, err
}

// Total implements domain.SearchToolRepository.
func (r *SearchToolRepository) Total(ctx context.Context, filter domain.SearchFilter) (int64, error) {
	panic("unimplemented")
}
