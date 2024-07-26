package repository

import (
	"context"
	"fmt"
	"go-service/internal/notification/domain"
	domain_notification "go-service/internal/notification/domain"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"
	"reflect"

	"gorm.io/gorm"
)

type notificationRepository struct {
	table      string
	buildParam func(n int) string
	db         *gorm.DB
	toArray    pq.Array
	modelType  reflect.Type
	logger     *logger.Logger
}

func NewNotificationRepository(table string, buildParam func(int) string, db *gorm.DB, logger *logger.Logger, toArray pq.Array) domain_notification.NotificationRepository {
	modelType := reflect.TypeOf(domain_notification.Notification{})
	return &notificationRepository{
		table:      table,
		buildParam: buildParam,
		db:         db,
		toArray:    toArray,
		modelType:  modelType,
		logger:     logger,
	}
}

func (r *notificationRepository) Search(ctx context.Context, filter domain.NotificationFilter) ([]domain.Notification, error) {
	list := []domain.Notification{}
	qr, params := r.buildQuery(filter)
	db := sql.GetTx(ctx, r.db)
	err := sql.QueryWithArray(db, &list, qr, r.toArray, params)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}
	return list, nil
}

func (r *notificationRepository) buildQuery(filter domain.NotificationFilter) (string, []interface{}) {
	qr := "select * from " + r.table

	params := []interface{}{}
	filterClause := "where"
	if filter.CreatedFrom != nil {
		params = append(params, filter.CreatedFrom)
		filterClause = filterClause + fmt.Sprintf("created_at >= %s", r.buildParam(len(params)))
	}

	if filter.CreatedTo != nil {
		params = append(params, filter.CreatedTo)
		filterClause = filterClause + fmt.Sprintf("created_at <= %s", r.buildParam(len(params)))
	}
	if filter.SubscriberId != nil {
		params = append(params, filter.SubscriberId)
		filterClause = filterClause + fmt.Sprintf(`Subscriber @> '{"subscriber_id": %s}'`, r.buildParam(len(params)))
	}
	filterClause = fmt.Sprintf("%s %s and 1 = 1", qr, filterClause)
	return filterClause, params
}

func (r *notificationRepository) Total(ctx context.Context, clientID string) (int64, error) {
	var total int64
	qr := "Select count(*) from %s where userId = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	err := sql.QueryWithArray(r.db, &total, stmt, r.toArray, clientID)
	return total, err
}

func (r *notificationRepository) TotalUnread(ctx context.Context, clientID string) (int64, error) {
	var total int64
	qr := "Select count(*) from %s where userId = %s and is_read = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2))
	err := sql.QueryWithArray(r.db, &total, stmt, r.toArray, clientID, false)
	return total, err
}

func (r *notificationRepository) Insert(ctx context.Context, notification domain_notification.Notification) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr, param, err := sql.BuildToInsert(db, r.table, notification, r.buildParam, r.modelType)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}

	res, err := sql.Exec(db, qr, param...)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return res, err
}
func (r *notificationRepository) Patch(ctx context.Context, notification map[string]interface{}) (int64, error) {
	panic("")
}
func (r *notificationRepository) Delete(ctx context.Context, id string) (int64, error) {
	panic("")
}
