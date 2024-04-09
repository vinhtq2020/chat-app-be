package repository

import (
	"fmt"
	"go-service/internal/search/domain"
	"go-service/pkg/database/sql"
	"go-service/pkg/database/sql/pq"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SearchRepository struct {
	table   string
	db      *gorm.DB
	toArray pq.Array
}

func NewSearchRepository(table string, db *gorm.DB, toArray pq.Array) *SearchRepository {
	return &SearchRepository{table: table, db: db, toArray: toArray}
}

func (f *SearchRepository) Search(e *gin.Context, result interface{}, filter domain.SearchFilter) error {
	qr := f.buildFilter(filter)
	err := sql.QueryWithArray(f.db, result, qr, f.toArray)
	return err
}

func (f *SearchRepository) Total(e *gin.Context) (int64, error) {
	total := int64(0)
	qr := fmt.Sprintf("select count(*) from %s", f.table)
	err := sql.Query(f.db, qr, &total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (f *SearchRepository) buildFilter(filter domain.SearchFilter) string {
	selectClause := fmt.Sprintf("select * from %s ", f.table)
	whereClause := "where"
	orderByClause := ""
	limitClause := ""
	if filter.Q != nil {
		whereClause += fmt.Sprintf("%s q like %%%s%% ", whereClause, *filter.Q)
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
	return selectClause + orderByClause + limitClause
}
