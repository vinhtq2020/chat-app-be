package repository

import (
	"fmt"
	"go-service/internal/user/domain"
	"go-service/pkg/database/sql"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository struct {
	db        *gorm.DB
	table     string
	modelType reflect.Type
}

func NewUserRepository(db *gorm.DB, table string) *UserRepository {
	modelType := reflect.TypeOf(domain.User{})
	return &UserRepository{db: db, modelType: modelType, table: table}
}

func (r *UserRepository) Load(e *gin.Context, id string) (domain.User, error) {
	var user domain.User
	qr := "Select * from users where id = $1"

	r.db.Raw(qr, "123").Scan(&user)
	return user, nil
}

func (r *UserRepository) buildParam(s int) string {
	return fmt.Sprintf("$%v", s)
}

func (r *UserRepository) Create(e *gin.Context, user domain.User) (int64, error) {
	qr, params, err := sql.BuildToInsert(r.db, r.table, user, r.buildParam, r.modelType)
	if err != nil {
		return -1, err
	}
	res, err := sql.Exec(r.db, qr, params...)
	return res, err
}
