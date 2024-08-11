package repository

import (
	"context"
	"fmt"
	"go-service/internal/friend/domain"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/logger"
	"reflect"

	"gorm.io/gorm"
)

type FriendRepository struct {
	table      string
	db         *gorm.DB
	buildParam func(int) string
	logger     *logger.Logger
	modelType  reflect.Type
}

func NewFriendRepository(db *gorm.DB, logger *logger.Logger, buildParam func(int) string) *FriendRepository {
	modelType := reflect.TypeOf(domain.Friend{})
	return &FriendRepository{
		table:      "friends",
		db:         db,
		logger:     logger,
		modelType:  modelType,
		buildParam: buildParam,
	}
}

func (r *FriendRepository) AlreadyFriended(ctx context.Context, userId1 string, userId2 string) (int64, error) {
	var res int64
	qr := "select count(*) from %s where (user_id1 = %s and user_id2 = %s) or (user_id1 = %s and user_id2 = %s)"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2), r.buildParam(2), r.buildParam(1))
	err := sql.Query(r.db, stmt, &res, r.logger, userId1, userId2)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return res, nil
}

func (r *FriendRepository) Upsert(ctx context.Context, friend domain.Friend) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr := "insert into %s(user_id1, user_id2, status) values (%s, %s, %s) on conflict(user_id1, user_id2) do update set status = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2), r.buildParam(3), r.buildParam(3))
	res, err := sql.Exec(db, stmt, r.logger, friend.UserId1, friend.UserId2, friend.Status)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return res, nil
}
func (r *FriendRepository) UpdateStatus(ctx context.Context, userId1 string, userId2 string, status string) (int64, error) {
	qr := "update %s where (user_id1 = %s and user_id2 = %s) or (user_id1 = %s and user_id = %s) set status = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2), r.buildParam(2), r.buildParam(1), r.buildParam(3))
	res, err := sql.Exec(r.db, stmt, r.logger, userId1, userId2, status)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return res, nil

}
func (r *FriendRepository) Delete(ctx context.Context, id string) (int64, error) {
	panic("")
}
