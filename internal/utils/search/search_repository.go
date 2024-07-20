package search

import (
	"context"
	"fmt"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"

	"gorm.io/gorm"
)

type searchRepository struct {
	table   string
	db      *gorm.DB
	toArray pq.Array
}

func NewSearchRepository(table string, db *gorm.DB, toArray pq.Array) *searchRepository {
	return &searchRepository{table: table, db: db, toArray: toArray}
}

func (f *searchRepository) Search(ctx context.Context, result interface{}, filter SearchFilter) error {
	params, qr := BuildQuery(f.table, filter)
	err := sql.QueryWithArray(f.db, result, qr, f.toArray, params...)
	return err
}

func (f *searchRepository) Total(ctx context.Context) (int64, error) {
	total := int64(0)
	qr := fmt.Sprintf("select count(*) from %s", f.table)
	err := sql.Query(f.db, qr, &total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func BuildQuery(table string, filter SearchFilter) (params []interface{}, qr string) {
	selectClause := fmt.Sprintf("select * from %s ", table)
	params, filterClause := BuildFilter(filter)
	return params, selectClause + filterClause
}

func BuildFilter(filter SearchFilter, opts ...interface{}) (params []interface{}, qr string) {
	buildParam := sql.BuildParam
	if len(opts) > 0 {
		buildParam = opts[0].(func(n int) string)
	}
	whereClause := "where"
	orderByClause := ""
	limitClause := ""
	if filter.Q != nil {
		params = append(params, *filter.Q)
		whereClause = fmt.Sprintf("%s q like CONCAT('%%',%s::text,'%%')", whereClause, buildParam(len(params)))
	}

	if len(filter.Sorts) > 0 {
		for _, v := range filter.Sorts {
			if len(v) > 0 {
				sortType := "ASC"
				if v[0] == '-' {
					sortType = "DESC"
				}
				orderByClause = fmt.Sprintf(" %s %s %s,", orderByClause, v, sortType)
			}

		}

		orderByClause = "order by" + orderByClause[:len(orderByClause)-2]
	}

	if filter.Page != nil && filter.Limit != nil {
		offset := *filter.Page * *filter.Limit
		limitClause = fmt.Sprintf("LIMIT %v OFFSET %v ", filter.Limit, offset)
	}
	return params, whereClause + orderByClause + limitClause
}
