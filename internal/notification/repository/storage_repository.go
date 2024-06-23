package repository

import (
	"context"
	"fmt"
	"go-service/internal/notification/domain"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
	"reflect"

	"gorm.io/gorm"
)

type notificationStorageRepository struct {
	table      string
	buildParam func(n int) string
	db         *gorm.DB
	toArray    pq.Array
	modelType  reflect.Type
}

func NewStorageRepository(table string, buildParam func(int) string, db *gorm.DB, toArray pq.Array) domain.NotificationStorageRepository {
	return &notificationStorageRepository{
		table:      table,
		buildParam: buildParam,
		db:         db,
		toArray:    toArray,
	}
}

func (r *notificationStorageRepository) Total(ctx context.Context, clientID string) (int64, error) {
	var total int64
	qr := "Select count(*) from %s where userId = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	err := sql.QueryWithArray(r.db, &total, stmt, r.toArray, clientID)
	return total, err
}

func (r *notificationStorageRepository) TotalUnread(ctx context.Context, clientID string) (int64, error) {
	var total int64
	qr := "Select count(*) from %s where userId = %s and is_read = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2))
	err := sql.QueryWithArray(r.db, &total, stmt, r.toArray, clientID, false)
	return total, err
}

func (r *notificationStorageRepository) Insert(ctx context.Context, notification domain.Notification) (int64, error) {
	qr, param, err := sql.BuildToInsert(r.db, r.table, notification, r.buildParam, r.modelType)
	if err != nil {
		return -1, err
	}

	res, err := sql.Exec(r.db, qr, param...)
	return res, err
}
func (r *notificationStorageRepository) Patch(ctx context.Context, notification map[string]interface{}) (int64, error) {
	panic("")
}
func (r *notificationStorageRepository) Delete(ctx context.Context, id string) (int64, error) {
	panic("")
}
