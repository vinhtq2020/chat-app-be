package repository

import (
	"fmt"
	"go-service/internal/sequence/domain"
	"go-service/pkg/sql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SequenceRepository struct {
	db          *gorm.DB
	table       string
	sequenceCol string
	nameCol     string
	buildParams func(int64) string
}

func NewSequenceRepository(db *gorm.DB, table string, options ...interface{}) *SequenceRepository {
	sequenceCol := "sequence"
	nameCol := "name"
	buildParams := func(n int64) string {
		return fmt.Sprintf("$%v", n)
	}
	if len(options) > 3 {
		sequenceCol = options[0].(string)
		nameCol = options[1].(string)
		buildParams = options[3].(func(n int64) string)
	}
	if len(options) > 2 {
		sequenceCol = options[0].(string)
		nameCol = options[1].(string)
	}

	if len(options) > 1 {
		sequenceCol = options[0].(string)
	}

	return &SequenceRepository{
		db:          db,
		table:       table,
		sequenceCol: sequenceCol,
		nameCol:     nameCol,
		buildParams: buildParams,
	}
}

func (r *SequenceRepository) Next(c *gin.Context, module string) (int64, error) {
	qr := fmt.Sprintf(
		"insert into %s as s values(%s, 1) on conflict(%s) do update set %s = s.%s + 1 where s.%s = %s",
		r.table, r.buildParams(1), r.nameCol, r.sequenceCol, r.sequenceCol, r.nameCol, r.buildParams(1))

	return sql.Exec(r.db, qr, module)
}

func (r *SequenceRepository) GetSequence(c *gin.Context, module string) (int64, error) {
	var sequence domain.Sequence
	qr := fmt.Sprintf("select * from %s where name = %s", r.table, r.buildParams(1))
	err := sql.Query(r.db, qr, &sequence, module)
	if err != nil {
		return -1, err
	}
	return sequence.SequenceNo, nil
}
