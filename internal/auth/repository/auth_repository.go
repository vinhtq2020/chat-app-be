package repository

import (
	"fmt"
	"go-service/internal/auth/domain"
	"go-service/pkg/sql"
	"go-service/pkg/sql/pq"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db        *gorm.DB
	table     string
	modelType reflect.Type
	toArray   pq.Array
}

func NewAuthRepository(db *gorm.DB, table string) *AuthRepository {
	modelType := reflect.TypeOf(domain.UserLoginData{})
	return &AuthRepository{db: db, modelType: modelType, table: table}
}

func (r *AuthRepository) Register(e *gin.Context, dt domain.UserLoginData) (int64, error) {
	qr, params, err := sql.BuildToInsert(r.db, r.table, dt, r.buildParam, r.modelType)
	if err != nil {
		return -1, err
	}
	res, err := sql.Exec(r.db, qr, params...)
	return res, err
}

func (r *AuthRepository) InTransaction(e *gin.Context, tx func() (int64, error)) (int64, error) {
	return sql.ExecuteTx(r.db, tx)
}

func (r *AuthRepository) GetUserLoginData(e *gin.Context, email string) (*domain.UserLoginData, error) {
	var users []domain.UserLoginData
	qr := "select * from %s where email = %s"
	stmt := fmt.Sprintf(qr, r.table, r.buildParam(1))
	err := sql.QueryWithArray(r.db, users, stmt, r.toArray, email)
	if err != nil || len(users) > 0 {
		return nil, err
	}
	return &users[0], nil
}

func (r *AuthRepository) buildParam(num int) string {
	return fmt.Sprintf("$%v", num)
}
