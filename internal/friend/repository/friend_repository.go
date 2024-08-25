package repository

import (
	"context"
	"fmt"
	"go-service/internal/friend/domain"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
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
	keys       []string
	toArray    pq.Array
}

func NewFriendRepository(db *gorm.DB, logger *logger.Logger, buildParam func(int) string, toArray pq.Array) *FriendRepository {
	modelType := reflect.TypeOf(domain.Relation{})
	keys := sql.GetPrimaryKeys(modelType)
	return &FriendRepository{
		table:      "relations",
		db:         db,
		logger:     logger,
		modelType:  modelType,
		buildParam: buildParam,
		keys:       keys,
		toArray:    toArray,
	}
}

func (r *FriendRepository) InTransaction(ctx context.Context, ex func(ctx context.Context, db *gorm.DB) (int64, error)) (int64, error) {
	return sql.ExecuteTx(ctx, r.db, ex)
}

func (r *FriendRepository) Load(ctx context.Context, userId1 string, userId2 string, relationType domain.RelationType) (*domain.Relation, error) {
	var res []domain.Relation
	db := sql.GetTx(ctx, r.db)
	qr := "select * from %s where ((user_id1 = %s and user_id2 = %s) or (user_id1 = %s and user_id2 = %s)) and relation_type = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2), r.buildParam(2), r.buildParam(1), r.buildParam(3))
	err := sql.QueryWithArray(db, &res, stmt, r.toArray, r.logger, userId1, userId2, relationType)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (r *FriendRepository) Upsert(ctx context.Context, friend domain.Relation) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr := `insert into %s(user_id1, user_id2, relation_type, status, created_at, created_by, updated_at, updated_by, notification_id)
		values (%s, %s, %s, %s, %s, %s, %s, %s, %s) on conflict (user_id1, user_id2, relation_type)
		do update set status = %s, updated_at = %s, updated_by = %s, notification_id = %s`
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2), r.buildParam(3), r.buildParam(4),
		r.buildParam(5), r.buildParam(6), r.buildParam(7), r.buildParam(8), r.buildParam(9),
		r.buildParam(4), r.buildParam(7), r.buildParam(8), r.buildParam(9))

	res, err := sql.Exec(db, stmt, r.logger, friend.UserId1, friend.UserId2, friend.RelationType, friend.Status, friend.CreatedAt, friend.CreatedBy, friend.UpdatedAt, friend.UpdatedBy, friend.NotificationId)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return res, nil
}
func (r *FriendRepository) Patch(ctx context.Context, friendRelation map[string]interface{}) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr, params, err := sql.BuildToPatch(r.db, r.table, r.modelType, friendRelation, r.keys, r.buildParam)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	res, err := sql.Exec(db, qr, r.logger, params...)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return res, nil

}
func (r *FriendRepository) Delete(ctx context.Context, id string) (int64, error) {
	panic("")
}
