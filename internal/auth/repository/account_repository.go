package repository

import (
	"context"
	"fmt"
	"go-service/internal/auth/domain"
	"go-service/pkg/database/postgres"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"
	"reflect"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db               *gorm.DB
	table            string
	modelType        reflect.Type
	refreshTokenType reflect.Type
	refreshTokenPk   []string
	logger           *logger.Logger
	toArray          pq.Array
}

func NewAccountRepository(db *gorm.DB, table string, logger *logger.Logger, toArray pq.Array) *AccountRepository {
	modelType := reflect.TypeOf(domain.Account{})
	refreshTokenType := reflect.TypeOf(domain.RefreshToken{})
	refreshTokenPk := postgres.GetPrimaryKeys(refreshTokenType)

	return &AccountRepository{db: db, table: table, modelType: modelType, refreshTokenType: refreshTokenType,
		logger: logger, toArray: toArray,
		refreshTokenPk: refreshTokenPk}
}

func (r *AccountRepository) Insert(ctx context.Context, dt domain.Account) (int64, error) {
	db := r.db

	tx, exist := ctx.Value("tx").(*gorm.DB)
	if exist {
		db = tx
	}

	qr, params, err := sql.BuildToInsert(db, r.table, dt, r.buildParam, r.modelType)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	_, err = sql.Exec(db, qr, params...)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return 1, nil
}

func (r *AccountRepository) InTransaction(ctx context.Context, ex func(db *gorm.DB) (int64, error)) (int64, error) {
	return sql.ExecuteTx(ctx, r.db, ex)
}

func (r *AccountRepository) Exist(ctx context.Context, email string) (int64, error) {
	qr := "select count(*) from %s where email = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	res := int64(0)
	err := sql.Query(r.db, stmt, &res, email)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}

	return res, nil
}

func (r *AccountRepository) Load(ctx context.Context, email string) (*domain.Account, error) {
	var users []domain.Account
	qr := "select * from %s where email = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	err := sql.QueryWithArray(r.db, &users, stmt, r.toArray, email)
	if err != nil || len(users) == 0 {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}
	return &users[0], nil
}

func (r *AccountRepository) buildParam(num int) string {
	return fmt.Sprintf("$%v", num)
}
