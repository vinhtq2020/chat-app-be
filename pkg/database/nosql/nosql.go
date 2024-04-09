package nosql

import (
	"errors"
	"go-service/pkg/database/cassandra"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func Query(ctx *gin.Context, db *cassandra.CassandraDBConnection, result interface{}, query string, value ...interface{}) error {
	iter := db.Session.Query(query).Bind(value...).WithContext(ctx).Iter()
	scanner := iter.Scanner()
	cols := iter.Columns()
	for scanner.Next() {
		fields, err := StructScan(result, cols)
		if err != nil {
			return err
		}

		err = scanner.Scan(fields...)
		if err != nil {
			return err
		}
	}
	return nil
}

func getTagChild(tags string, tagChild string, options ...bool) string {
	haveSemicolon := true
	if len(options) > 0 {
		haveSemicolon = options[0]
	}
	parts := strings.Split(tags, ";")
	for _, part := range parts {
		if strings.Contains(part, tagChild) {
			if haveSemicolon {
				return strings.Split(part, ":")[1]
			}
			return part
		}
	}
	return ""
}

func GetFieldIndexes(modelType reflect.Type) (map[string]int, error) {
	ci := make(map[string]int, 0)
	if modelType.Kind() != reflect.Struct {
		return ci, errors.New("type is not struct")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagName := field.Tag.Get("gorm")
		if len(tagName) > 0 {
			tagChild := getTagChild(tagName, "column")
			ci[tagChild] = i
		}
	}
	return ci, nil
}

func StructScan(src interface{}, columns []gocql.ColumnInfo) (r []interface{}, err error) {
	modelType := reflect.TypeOf(src).Elem()
	modelValue := reflect.Indirect(reflect.ValueOf(src))
	fieldIndexes, err := GetFieldIndexes(modelType)

	if err != nil {
		return nil, err
	}

	for _, col := range columns {
		index, ok := fieldIndexes[col.Name]
		if !ok {
			var t interface{}
			r = append(r, t)
			continue
		}
		valueField := modelValue.Field(index)
		// typeField := modelType.Field(index)
		valueFieldTmp := valueField.Addr().Interface()
		// if typeField.Type.Kind() == reflect.Slice || typeField.Type.Kind() == reflect.Pointer && typeField.Type.Elem().Kind() == reflect.Slice {
		// 	valueFieldTmp = toArray(valueFieldTmp)
		// }
		r = append(r, valueFieldTmp)

	}
	return r, err
}
