package repository

import (
	"errors"
	"fmt"
	"go-service/internal/auth/domain"
	"go-service/pkg/database/sql"
	"go-service/pkg/database/sql/pq"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db               *gorm.DB
	table            string
	modelType        reflect.Type
	refreshTokenType reflect.Type
	toArray          pq.Array
}

func NewAuthRepository(db *gorm.DB, table string, toArray pq.Array) *AuthRepository {
	modelType := reflect.TypeOf(domain.UserLoginData{})
	refreshTokenType := reflect.TypeOf(domain.RefreshToken{})
	return &AuthRepository{db: db, table: table, modelType: modelType, refreshTokenType: refreshTokenType, toArray: toArray}
}

func (r *AuthRepository) Register(e *gin.Context, dt domain.UserLoginData) (int64, error) {
	db := r.db

	tx, exist := e.Get("tx")
	if exist {
		db = tx.(*gorm.DB)
	}
	qr, params, err := sql.BuildToInsert(db, r.table, dt, r.buildParam, r.modelType)
	if err != nil {
		return -1, err
	}
	res, err := sql.Exec(db, qr, params...)
	return res, err
}

func (r *AuthRepository) InTransaction(e *gin.Context, ex func(db *gorm.DB) (int64, error)) (int64, error) {
	return sql.ExecuteTx(e, r.db, ex)
}

func (r *AuthRepository) Exist(e *gin.Context, email string) (int64, error) {
	qr := "select count(*) from %s where email = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	res := int64(0)
	err := sql.Query(r.db, stmt, &res, email)
	if err != nil {
		return -1, err
	}

	return res, nil
}

func (r *AuthRepository) GetUserLoginData(e *gin.Context, email string) (*domain.UserLoginData, error) {
	var users []domain.UserLoginData
	qr := "select * from %s where email = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	err := sql.QueryWithArray(r.db, &users, stmt, r.toArray, email)
	if err != nil || len(users) == 0 {
		return nil, err
	}
	return &users[0], nil
}

func (r *AuthRepository) AddRefreshToken(e *gin.Context, refreshToken domain.RefreshToken) (int64, error) {
	qr, params, err := sql.BuildToInsert(r.db, "refresh_tokens", refreshToken, r.buildParam, r.refreshTokenType)
	if err != nil {
		return -1, err
	}
	res, err := sql.Exec(r.db, qr, params...)
	return res, err
}

func (r *AuthRepository) RemoveRefreshToken(e *gin.Context, userId string, ipAddress string, deviceId string, browser string) (int64, error) {
	qr := "delete from refresh_tokens where user_id = $1 and ip_address = $2 and device_id = $3 and browser = $4"
	res, err := sql.Exec(r.db, qr, userId, ipAddress, deviceId, browser)
	if err != nil {
		return -1, err
	}
	if res == 0 {
		return 0, errors.New("row not found")
	}
	return res, nil
}
func (r *AuthRepository) buildParam(num int) string {
	return fmt.Sprintf("$%v", num)
}
