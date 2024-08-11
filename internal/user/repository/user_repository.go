package repository

import (
	"context"
	"fmt"
	user_domain "go-service/internal/user/domain"
	"go-service/pkg/database/postgres"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"
	"reflect"

	"gorm.io/gorm"
)

type UserRepository struct {
	db        *gorm.DB
	table     string
	logger    *logger.Logger
	modelType reflect.Type
	toArray   pq.Array
}

func NewUserRepository(db *gorm.DB, table string, logger *logger.Logger, toArray pq.Array) *UserRepository {
	modelType := reflect.TypeOf(user_domain.User{})
	return &UserRepository{db: db, modelType: modelType, table: table, toArray: toArray}
}

func (r *UserRepository) Load(ctx context.Context, id string) (*user_domain.User, error) {
	var res []user_domain.User
	stmt := fmt.Sprintf("select * from %s where id = %s", r.table, r.buildParam(1))
	err := postgres.QueryWithArray(r.db, &res, stmt, r.toArray, r.logger, id)
	if err != nil || len(res) == 0 {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}
	return &res[0], nil
}

func (r *UserRepository) buildParam(s int) string {
	return fmt.Sprintf("$%v", s)
}

func (r *UserRepository) Create(ctx context.Context, user user_domain.User) (int64, error) {
	qr, params, err := sql.BuildToInsert(r.db, r.table, user, r.buildParam, r.modelType)
	if err != nil {
		return -1, err
	}
	res, err := sql.Exec(r.db, qr, r.logger, params...)
	return res, err
}

func (r *UserRepository) Exist(ctx context.Context, id string) (bool, error) {
	qr := "select count(*) from %s where id = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	res := int64(0)
	err := sql.Query(r.db, stmt, &res, r.logger, id)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return false, err
	}

	return res > 0, nil
}
