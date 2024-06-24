package repository

import (
	"context"
	"fmt"
	"go-service/internal/user/domain"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/logger"
	"reflect"

	"gorm.io/gorm"
)

type UserRepository struct {
	db        *gorm.DB
	table     string
	logger    *logger.Logger
	modelType reflect.Type
}

func NewUserRepository(db *gorm.DB, table string, logger *logger.Logger) *UserRepository {
	modelType := reflect.TypeOf(domain.User{})
	return &UserRepository{db: db, modelType: modelType, table: table}
}

func (r *UserRepository) Load(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	qr := "Select * from users where id = $1"

	r.db.Raw(qr, "123").Scan(&user)
	return user, nil
}

func (r *UserRepository) buildParam(s int) string {
	return fmt.Sprintf("$%v", s)
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	qr, params, err := sql.BuildToInsert(r.db, r.table, user, r.buildParam, r.modelType)
	if err != nil {
		return -1, err
	}
	res, err := sql.Exec(r.db, qr, params...)
	return res, err
}

func (r *UserRepository) Exist(ctx context.Context, id string) (bool, error) {
	qr := "select count(*) from %s where id = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	res := int64(0)
	err := sql.Query(r.db, stmt, &res, id)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return false, err
	}

	return res > 0, nil
}
