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

type RefreshTokenRepository struct {
	db          *gorm.DB
	table       string
	modelType   reflect.Type
	primaryKeys []string

	logger  *logger.Logger
	toArray pq.Array
}

func NewRefreshTokenRepository(db *gorm.DB, table string, logger *logger.Logger, toArray pq.Array) *RefreshTokenRepository {
	modelType := reflect.TypeOf(domain.RefreshToken{})
	primaryKeys := postgres.GetPrimaryKeys(modelType)
	return &RefreshTokenRepository{
		db:          db,
		table:       table,
		modelType:   modelType,
		primaryKeys: primaryKeys,
		logger:      logger,
		toArray:     toArray,
	}
}

func (r *RefreshTokenRepository) InTransaction(ctx context.Context, ex func(db *gorm.DB) (int64, error)) (int64, error) {
	return sql.ExecuteTx(ctx, r.db, ex)
}

func (r *RefreshTokenRepository) Load(ctx context.Context, browser string, ip string, deviceId string) (*domain.RefreshToken, error) {
	var res []domain.RefreshToken
	db := sql.GetTx(ctx, r.db)
	qr := fmt.Sprintf("select * from %s where browser = %s and ip_address = %s and device_id = %s", r.table, r.buildParam(1), r.buildParam(2), r.buildParam(3))
	err := sql.QueryWithArray(db, &res, qr, r.toArray, browser, ip, deviceId)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res[0], nil
}

func (r *RefreshTokenRepository) Insert(ctx context.Context, refreshToken domain.RefreshToken) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr, params, err := sql.BuildToInsert(db, r.table, refreshToken, r.buildParam, r.modelType)
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

func (r *RefreshTokenRepository) Patch(ctx context.Context, refreshToken map[string]interface{}) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr, params, err := sql.BuildToPatch(db, r.table, refreshToken, r.primaryKeys, r.buildParam)
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

func (r *RefreshTokenRepository) Delete(ctx context.Context, userId string, ipAddress string, deviceId string, browser string) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr := fmt.Sprintf("delete from %s where user_id = $1 and ip_address = $2 and device_id = $3 and browser = $4", r.table)
	res, err := sql.Exec(db, qr, userId, ipAddress, deviceId, browser)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}

	return res, nil
}

func (r *RefreshTokenRepository) buildParam(num int) string {
	return fmt.Sprintf("$%v", num)
}
