package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"go-service/pkg/database/postgres/pq"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func Query(db *gorm.DB, sql string, result interface{}, value ...interface{}) error {
	err := db.Raw(sql, value...).Scan(&result).Error
	return err
}

func QueryWithArray(db *gorm.DB, results interface{}, sql string, toArray func(a interface{}) interface {
	driver.Valuer
	sql.Scanner
}, value ...interface{}) error {
	// pointer => array => item
	modelType := reflect.TypeOf(results).Elem().Elem()

	tx := db.Raw(sql, value...)
	err := tx.Error
	if err != nil {
		return err
	}
	rows, err := tx.Rows()
	if err != nil {
		return err
	}
	colIndexes, err := rows.Columns()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {

		item := reflect.New(modelType).Interface()
		fields, err := StructScan(item, colIndexes, toArray)
		if err != nil {
			return err
		}
		err = rows.Scan(fields...)
		if err != nil {
			return err
		}

		ArrayAppend(results, item)

	}

	return nil
}

func ArrayAppend(array interface{}, item interface{}) interface{} {
	arrayValue := reflect.ValueOf(array)
	arrayValueElem := reflect.Indirect(arrayValue)
	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Pointer {
		itemValue = reflect.Indirect(itemValue)
	}
	arrayValueElem.Set(reflect.Append(arrayValueElem, itemValue))
	return array
}

func Exec(db *gorm.DB, sql string, value ...interface{}) (int64, error) {
	tx := db.Exec(sql, value...)
	return tx.RowsAffected, tx.Error
}

func buildParams(n int, buildParam func(int) string) []string {
	params := []string{}
	for i := 1; i <= n; i++ {
		params = append(params, buildParam(i))
	}
	return params
}

func BuildParam(n int) string {
	return fmt.Sprintf("$%d", n)
}

func GetPrimaryKeys(modelType reflect.Type) (keys []string) {
	// return GetTags(modelType, "gorm", "primaryKey", false)
	for i := 0; i < modelType.NumField(); i++ {
		v, ok := modelType.Field(i).Tag.Lookup("gorm")
		if ok {
			vs := strings.Split(v, ";")
			hasKey := false
			column := ""
			for i := 0; i < len(vs); i++ {
				if vs[i] == "primaryKey" {
					hasKey = true
				} else if strings.Contains(vs[i], "column") {
					t := strings.Split(vs[i], ":")
					if len(t) > 1 {
						column = t[1]
					}
				}
			}
			if hasKey && len(column) > 0 {
				keys = append(keys, column)
			}
		}
	}
	return keys
}

func GetValuesIfTagExist(obj interface{}, tag string) []interface{} {
	params := []interface{}{}
	modelValue := reflect.ValueOf(obj)
	for i := 0; i < modelValue.NumField(); i++ {
		fieldValue := modelValue.Field(i).Interface()
		fieldTag := modelValue.Type().Field(i).Tag.Get(tag)
		if len(fieldTag) > 0 {
			params = append(params, fieldValue)
		}
	}
	return params
}

func BuildToInsert(db *gorm.DB, table string, obj interface{}, buildParam func(int) string, modelType reflect.Type) (string, []interface{}, error) {
	tags := GetTags(modelType, "gorm", "column")
	if len(tags) == 0 {
		return "", nil, errors.New("missing gorm's tag in struct field")
	}
	qr := "insert into %s(%s) values(%s)"

	stmt := fmt.Sprintf(qr, table, strings.Join(tags, ","), strings.Join(buildParams(len(tags), buildParam), ", "))
	params := GetValuesIfTagExist(obj, "gorm")
	return stmt, params, nil
}

func Exists(arr []string, item string) bool {
	for _, v := range arr {
		if item == v {
			return true
		}
	}
	return false
}

func BuildToPatch(db *gorm.DB, table string, params map[string]interface{}, keys []string, buildParam func(int) string) (string, []interface{}, error) {
	set := []string{}
	setValue := []interface{}{}
	i := 1
	where := []string{}
	for k, v := range params {
		if Exists(keys, k) {
			setValue = append(setValue, v)
			where = append(where, fmt.Sprintf("%v=%v", k, buildParam(i)))
			i++
			continue
		}
		set = append(set, fmt.Sprintf("%v=%v", k, buildParam(i)))
		setValue = append(setValue, v)
		i++
	}
	if len(where) != len(keys) {
		return "", nil, errors.New("not have full primary keys")
	}
	sql := fmt.Sprintf("update %v set %v where %v", table, strings.Join(set, ", "), strings.Join(where, "and "))
	return sql, setValue, nil
}

func GetTags(modelType reflect.Type, tag string, options ...interface{}) []string {
	tagChild := ""
	semicolon := true
	if len(options) > 1 {
		tagChild = options[0].(string)
		semicolon = options[1].(bool)
	} else if len(options) > 0 {
		tagChild = options[0].(string)
	}
	tags := []string{}

	if modelType.Kind() == reflect.Slice {
		modelType = modelType.Elem()
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagName := field.Tag.Get(tag)

		if len(tagName) > 0 {
			if len(tagChild) <= 0 {
				tags = append(tags, tagName)
				continue
			}
			tc := getTagChild(tagName, tagChild, semicolon)
			tags = append(tags, tc)
		}
	}
	return tags
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

func StructScan(src interface{}, columnIndexes []string, toArray pq.Array) (r []interface{}, err error) {
	modelType := reflect.TypeOf(src).Elem()
	modelValue := reflect.Indirect(reflect.ValueOf(src))
	fieldIndexes, err := GetFieldIndexes(modelType)

	if err != nil {
		return nil, err
	}

	for _, colName := range columnIndexes {
		index, ok := fieldIndexes[colName]
		if !ok {
			var t interface{}
			r = append(r, t)
			continue
		}
		valueField := modelValue.Field(index)
		typeField := modelType.Field(index)
		valueFieldTmp := valueField.Addr().Interface()
		if typeField.Type.Kind() == reflect.Slice || typeField.Type.Kind() == reflect.Pointer && typeField.Type.Elem().Kind() == reflect.Slice {
			valueFieldTmp = toArray(valueFieldTmp)
		}
		r = append(r, valueFieldTmp)

	}
	return r, err
}

func ExecuteTx(ctx context.Context, db *gorm.DB, ex func(tx *gorm.DB) (int64, error)) (int64, error) {
	tx := db.Begin()
	ctx = context.WithValue(ctx, "tx", tx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return -1, err
	}

	res, err := ex(tx)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	return res, tx.Commit().Error
}

func GetTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, exist := ctx.Value("tx").(*gorm.DB)
	if exist {
		db = tx
	}

	return db
}
