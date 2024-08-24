package repository

import (
	"context"
	"fmt"
	"go-service/internal/notification/domain"
	sql "go-service/pkg/database/postgres"
	"go-service/pkg/database/postgres/pq"
	"go-service/pkg/logger"
	"reflect"

	"gorm.io/gorm"
)

type notificationRepository struct {
	keys       []string
	table      string
	buildParam func(n int) string
	db         *gorm.DB
	toArray    pq.Array
	modelType  reflect.Type
	logger     *logger.Logger
}

func NewNotificationRepository(table string, buildParam func(int) string, db *gorm.DB, logger *logger.Logger, toArray pq.Array) domain.NotificationRepository {
	modelType := reflect.TypeOf(domain.Notification{})
	keys := sql.GetPrimaryKeys(modelType)
	return &notificationRepository{
		table:      table,
		buildParam: buildParam,
		db:         db,
		toArray:    toArray,
		modelType:  modelType,
		logger:     logger,
		keys:       keys,
	}
}

func (r *notificationRepository) Search(ctx context.Context, filter domain.NotificationFilter) ([]domain.Notification, error) {
	list := []domain.Notification{}
	qr, params := r.buildQuery(filter)
	db := sql.GetTx(ctx, r.db)
	err := sql.QueryWithArray(db, &list, qr, r.toArray, r.logger, params...)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}
	return list, nil
}

func (r *notificationRepository) buildQuery(filter domain.NotificationFilter) (string, []interface{}) {
	qr := "select * from " + r.table
	params := []interface{}{}
	filterClause := "where "
	if filter.CreatedFrom != nil {
		params = append(params, *filter.CreatedFrom)
		filterClause = filterClause + fmt.Sprintf("created_at >= %s and ", r.buildParam(len(params)))
	}

	if filter.CreatedTo != nil {
		params = append(params, *filter.CreatedTo)
		filterClause = filterClause + fmt.Sprintf("created_at <= %s and ", r.buildParam(len(params)))
	}
	if filter.SubscriberId != nil {
		params = append(params, *filter.SubscriberId)
		filterClause = filterClause + fmt.Sprintf(`EXISTS (
            SELECT 1
            FROM unnest(subscribers) AS json_obj
            WHERE json_obj->>'id' = %s
        ) and `, r.buildParam(len(params)))
	}
	if filter.Visible != nil {
		params = append(params, filter.Visible)
		filterClause = filterClause + fmt.Sprintf(`visible = %s' and`, r.buildParam(len(params)))
	}
	filterClause = fmt.Sprintf("%s %s 1 = 1", qr, filterClause)
	return filterClause, params
}

func (r *notificationRepository) Total(ctx context.Context, clientID string) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	var total int64
	qr := "Select count(*) from %s where userId = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	err := sql.QueryWithArray(db, &total, stmt, r.toArray, r.logger, clientID)
	return total, err
}

func (r *notificationRepository) TotalUnread(ctx context.Context, clientID string) (int64, error) {
	var total int64
	qr := "Select count(*) from %s where userId = %s and exist (select 1 from  unnest(subscribers) as b where b ->> 'isRead' = %s) "
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1), r.buildParam(2))
	err := sql.QueryWithArray(r.db, &total, stmt, r.toArray, r.logger, clientID, false)
	return total, err
}

func (r *notificationRepository) Insert(ctx context.Context, notification domain.Notification) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr, param, err := sql.BuildToInsert(db, r.table, notification, r.buildParam, r.modelType)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}

	res, err := sql.Exec(db, qr, r.logger, param...)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return -1, err
	}
	return res, err
}

func (r *notificationRepository) Load(ctx context.Context, notificationId string) (*domain.Notification, error) {
	var res []domain.Notification
	db := sql.GetTx(ctx, r.db)
	qr := "select * from %s where id = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	err := sql.QueryWithArray(db, &res, stmt, r.toArray, r.logger, notificationId)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
		return nil, err
	}
	return &res[0], nil
}

func (r *notificationRepository) Patch(ctx context.Context, notification map[string]interface{}) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr, params, err := sql.BuildToPatch(db, r.table, r.modelType, notification, r.keys, r.buildParam)
	if err != nil {
		return -1, err
	}
	res, err := sql.Exec(r.db, qr, r.logger, params...)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (r *notificationRepository) Delete(ctx context.Context, id string) (int64, error) {
	db := sql.GetTx(ctx, r.db)
	qr := "delete from %s where id = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	res, err := sql.Exec(db, stmt, r.logger, id)
	if err != nil {
		r.logger.LogError(err.Error(), nil)
	}
	return res, err
}
