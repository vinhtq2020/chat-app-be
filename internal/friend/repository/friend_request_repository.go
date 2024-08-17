package repository

import (
	"context"
	"fmt"
	"go-service/internal/friend/domain"
	friend_domain "go-service/internal/friend/domain"
	"go-service/pkg/database/postgres"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/logger"
	"reflect"

	"gorm.io/gorm"
)

type RequestFriendRepository struct {
	db         *gorm.DB
	logger     *logger.Logger
	buildParam func(int) string
	table      string
	userTable  string
	keys       []string
	modelType  reflect.Type
}

func NewRequestFriendRepository(db *gorm.DB, table string, user_table string, logger *logger.Logger, buidParam func(int) string) *RequestFriendRepository {
	modelType := reflect.TypeOf(friend_domain.FriendRequest{})
	primaryKeys := postgres.GetPrimaryKeys(modelType)

	return &RequestFriendRepository{
		db:         db,
		table:      table,
		logger:     logger,
		buildParam: buidParam,
		userTable:  user_table,
		keys:       primaryKeys,
		modelType:  modelType,
	}
}

func (r *RequestFriendRepository) Load(ctx context.Context, requestId string) (*domain.FriendRequest, error) {
	var frq domain.FriendRequest
	stmt := fmt.Sprintf("select * from %s where id = %s", r.table, r.buildParam(1))
	err := sql.Query(r.db, stmt, &frq, nil, requestId)
	if err != nil {
		return nil, err
	}
	return &frq, nil
}

func (r *RequestFriendRepository) Exist(ctx context.Context, requesterId string, requesteeId string, status domain.FriendRequestStatus) (bool, error) {
	qr := "select count(*) from %s where requester_id = %s and requestee_id = %s and status = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2), r.buildParam(3))
	var res int64
	err := sql.Query(r.db, stmt, &res, r.logger, requesterId, requesteeId, status)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return false, err
	}
	return res > 0, nil
}

func (r *RequestFriendRepository) All(ctx context.Context, userId string) ([]friend_domain.FriendRequest, error) {
	qr := "select * from %s where requester = %s and status = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2))
	var res []friend_domain.FriendRequest
	err := sql.Query(r.db, stmt, &res, r.logger, userId, friend_domain.StatusPending.Value())
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}
	return res, nil
}

func (r *RequestFriendRepository) Total(ctx context.Context) (int64, error) {
	qr := "select count(*) from %s"
	stmt := fmt.Sprintf(qr, r.table)
	var res int64
	err := sql.Query(r.db, stmt, &res, r.logger)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return res, nil
}
func (r *RequestFriendRepository) Create(ctx context.Context, friendRq friend_domain.FriendRequest) (int64, error) {
	db := postgres.GetTx(ctx, r.db)
	qr, params, err := sql.BuildToInsert(db, r.table, friendRq, r.buildParam, r.modelType)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}

	res, err := postgres.Exec(db, qr, r.logger, params...)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return res, nil
}
func (r *RequestFriendRepository) Patch(ctx context.Context, friendRq map[string]interface{}) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr, params, err := sql.BuildToPatch(r.db, r.table, r.modelType, friendRq, r.keys, r.buildParam)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	_, err = sql.Exec(db, qr, r.logger, params...)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return 1, nil
}
func (r *RequestFriendRepository) Delete(ctx context.Context, id string) (int64, error) {
	panic("")
}

func (r *RequestFriendRepository) InTransaction(ctx context.Context, ex func(ctx context.Context, db *gorm.DB) (int64, error)) (int64, error) {
	return sql.ExecuteTx(ctx, r.db, ex)
}
