package repository_test

import (
	"database/sql"
	"database/sql/driver"
	"go-service/internal/search/repository"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewSearchRepository(t *testing.T) {

	// mock gorm DB
	mockDb, _, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	db, _ := gorm.Open(dialector, &gorm.Config{})
	type Input struct {
		DB        *gorm.DB
		table     string
		modelType reflect.Type
		toArray   func(a interface{}) interface {
			driver.Valuer
			sql.Scanner
		}
	}
	type TestCase struct {
		Name     string
		Input    Input
		Expected bool
	}

	testCases := []TestCase{
		{
			Name: "Success Case",
			Input: Input{
				DB:        db,
				table:     "",
				modelType: reflect.TypeOf(10),
				toArray:   pq.Array,
			},
			Expected: true,
		},
	}

	for _, v := range testCases {
		t.Run(v.Name, func(t *testing.T) {
			repo := repository.NewSearchRepository(v.Input.table, v.Input.DB, v.Input.toArray)
			if repo == nil {
				t.Errorf("actual value: %v, expected %v", repo != nil, v.Expected)
			} else {
				t.Logf("actual value: %v, expected %v", repo != nil, v.Expected)
			}
		})
	}

}
func TestSearch(t *testing.T) {

}
