package search

import (
	"context"
	"fmt"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
	"reflect"

	"gorm.io/gorm"
)

type searchRepository[F any] struct {
	table      string
	db         *gorm.DB
	toArray    pq.Array
	buildQuery func(filter F, buildParam func(int) string) (string, []interface{})
	buildParam func(int) string
}

func NewSearchRepository[F any](table string, db *gorm.DB, toArray pq.Array, buildQuery func(filter F, buildParam func(int) string) (string, []interface{})) *searchRepository[F] {
	return &searchRepository[F]{table: table, db: db, toArray: toArray, buildQuery: buildQuery}
}

func (f *searchRepository[F]) Search(ctx context.Context, result interface{}, filter F) error {
	var params []interface{}
	var qr string
	if f.buildQuery == nil {
		qr, params = BuildQuery(f.table, filter)
	} else {
		qr, params = f.buildQuery(filter, f.buildParam)
	}
	err := sql.QueryWithArray(f.db, result, qr, f.toArray, nil, params...)
	return err
}

func (f *searchRepository[F]) Total(ctx context.Context) (int64, error) {
	total := int64(0)
	qr := fmt.Sprintf("select count(*) from %s", f.table)
	err := sql.Query(f.db, qr, &total, nil)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func BuildQuery(table string, filter interface{}) (qr string, params []interface{}) {
	selectClause := fmt.Sprintf("select * from %s ", table)
	filterClause, params := BuildFilter(filter)
	return selectClause + filterClause, params
}

func BuildFilter(filter interface{}, opts ...interface{}) (qr string, params []interface{}) {
	baseFilterType := reflect.TypeOf(&SearchFilter{})
	modelType := reflect.TypeOf(filter)
	modelVal := reflect.ValueOf(filter)
	var baseFilter SearchFilter
	var ok bool
	for i := 0; i < modelType.NumField(); i++ {
		if modelType.Field(i).Type == baseFilterType {
			baseFilter, ok = modelVal.Field(i).Interface().(SearchFilter)
			if !ok {
				return qr, params
			}
			break
		}
	}

	buildParam := sql.BuildParam
	if len(opts) > 0 {
		buildParam = opts[0].(func(n int) string)
	}
	whereClause := "where"
	orderByClause := ""
	limitClause := ""
	if baseFilter.Q != nil {
		params = append(params, *baseFilter.Q)
		whereClause = fmt.Sprintf("%s q like CONCAT('%%',%s::text,'%%')", whereClause, buildParam(len(params)))
	}

	if len(baseFilter.Sorts) > 0 {
		for _, v := range baseFilter.Sorts {
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

	if baseFilter.Page != nil && baseFilter.Limit != nil {
		offset := *baseFilter.Page * *baseFilter.Limit
		limitClause = fmt.Sprintf("LIMIT %v OFFSET %v ", baseFilter.Limit, offset)
	}
	return whereClause + orderByClause + limitClause, params
}
