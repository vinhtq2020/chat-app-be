package repository

import (
	"context"
	"fmt"
	friend_domain "go-service/internal/friend/domain"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/logger"

	"gorm.io/gorm"
)

type RequestFriendRepository struct {
	db         *gorm.DB
	logger     *logger.Logger
	buildParam func(int) string
	table      string
}

func NewRequestFriendRepository(db *gorm.DB, table string, logger *logger.Logger, buidParam func(int) string) *RequestFriendRepository {
	return &RequestFriendRepository{
		db:         db,
		table:      table,
		logger:     logger,
		buildParam: buidParam,
	}
}
func (r *RequestFriendRepository) Exist(ctx context.Context, userId string, friendId string) (bool, error) {
	qr := "select count(*) from %s where uid1 = %s and uid2 = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2))
	var res int64
	err := sql.Query(r.db, stmt, &res, userId, friendId)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return false, err
	}
	return res > 0, nil
}

func (r *RequestFriendRepository) UserExist(ctx context.Context, id string) (bool, error) {
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

func (r *RequestFriendRepository) All(ctx context.Context, userId string) ([]friend_domain.FriendRequest, error) {
	qr := "select * from %s  where uid1 = %s and status = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2))
	var res []friend_domain.FriendRequest
	err := sql.Query(r.db, stmt, &res, userId, friend_domain.StatusPending.Value())
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}
	return res, nil
}

func (r *RequestFriendRepository) Create(ctx context.Context, friendRq friend_domain.FriendRequest) (int64, error) {
	panic("")
}
func (r *RequestFriendRepository) Patch(ctx context.Context, friendRq map[string]interface{}) (int64, error) {
	panic("")
}
func (r *RequestFriendRepository) Delete(ctx context.Context, id string) (int64, error) {
	panic("")
}
