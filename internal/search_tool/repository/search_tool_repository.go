package repository

import (
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

func (r *SearchToolRepository) Search(filter domain.SearchFilter) ([]domain.SearchItem, error) {
	var res []domain.SearchItem
	selectClause := "select a.user_name, a.avatar_url from users a where a.user_name like CONCAT('%%',%s::text,'%%')"
	filterClause := search.BuildFilter(filter.SearchFilter)
	stmt := fmt.Sprintf(selectClause+filterClause, sql.BuildParam(1))
	err := sql.QueryWithArray(r.DB, &res, stmt, r.toArray, filter.Q)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
	}
	return res, err
}

// Total implements domain.SearchToolRepository.
func (r *SearchToolRepository) Total(filter domain.SearchFilter) (int64, error) {
	panic("unimplemented")
}
