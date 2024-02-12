package repository

import (
	"fmt"
	"go-service/internal/room/domain"
	"go-service/pkg/sql"
	"go-service/pkg/sql/pq"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db          *gorm.DB
	table       string
	toArray     pq.Array
	modelType   reflect.Type
	primaryKeys []string
}

func NewRoomReposiory(db *gorm.DB, table string, toArray pq.Array) *RoomRepository {
	modelType := reflect.TypeOf(domain.Room{})
	primaryKey := sql.GetPrimaryKeys(modelType)
	return &RoomRepository{
		db:          db,
		table:       table,
		toArray:     toArray,
		modelType:   modelType,
		primaryKeys: primaryKey,
	}
}

func (r *RoomRepository) buildParam(num int) string {
	return fmt.Sprintf("$%v", num)
}

func (r *RoomRepository) All(ctx *gin.Context) ([]domain.Room, error) {
	qr := fmt.Sprintf("select * from %s", r.table)
	var res []domain.Room
	// rows, err := r.db.Raw(qr).Rows()
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var item domain.Room
	// 	err := rows.Scan(&item.Id, &item.Name, r.toArray(&item.Members), &item.CreatedAt, &item.CreatedBy, &item.UpdatedAt, &item.UpdatedBy, &item.Version)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	res = append(res, item)
	// }
	err := sql.QueryWithArray(r.db, &res, qr, r.toArray)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Create implements domain.RoomRepository.
func (r *RoomRepository) Create(ctx *gin.Context, room domain.Room) (int64, error) {
	qr, param, err := sql.BuildToInsert(r.db, r.table, room, r.buildParam, r.modelType)
	if err != nil {
		return -1, err
	}
	res, err := sql.Exec(r.db, qr, param...)
	if err != nil {
		return -1, err
	}
	return res, nil
}

// Delete implements domain.RoomRepository.
func (r *RoomRepository) Delete(ctx *gin.Context, id string) (int64, error) {
	qr := "Delete from %s where id = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	res, err := sql.Exec(r.db, stmt, id)
	return res, err
}

// Load implements domain.RoomRepository.
func (r *RoomRepository) Load(ctx *gin.Context, id string) (*domain.Room, error) {
	var res []domain.Room
	qr := "select * from %s where id = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	err := sql.QueryWithArray(r.db, &res, stmt, r.toArray, id)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return &res[0], nil
}

// Patch implements domain.RoomRepository.
func (r *RoomRepository) Patch(ctx *gin.Context, room map[string]interface{}) (int64, error) {
	qr, vals, err := sql.BuildToPatch(r.db, r.table, room, r.primaryKeys, r.buildParam)
	if err != nil {
		return -1, err
	}
	res, err := sql.Exec(r.db, qr, vals...)
	if err != nil {
		return -1, err
	}
	return res, nil
}

// Update implements domain.RoomRepository.
func (*RoomRepository) Update(ctx *gin.Context, Room domain.Room) {
	panic("unimplemented")
}
