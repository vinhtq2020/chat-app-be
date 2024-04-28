package convert

import (
	"encoding/json"
	"reflect"
	"strings"
)

func ConvertObjectToMap(obj interface{}) map[string]interface{} {
	inrec, _ := json.Marshal(obj)
	var mapObj map[string]interface{}
	json.Unmarshal(inrec, &mapObj)
	return mapObj
}

// convert object to map. Field name of map is used tag json
func ToMapOmitEmpty(obj interface{}) map[string]interface{} {
	modelType := reflect.TypeOf(obj)
	modelValue := reflect.ValueOf(obj)
	mapRes := map[string]interface{}{}
	for i := 0; i < modelType.NumField(); i++ {
		tag, ok := modelType.Field(i).Tag.Lookup("json")
		if ok {
			tagChild := strings.Split(tag, ",")
			if !modelValue.Field(i).IsZero() {
				if modelValue.Field(i).Kind() == reflect.Pointer {
					mapRes[tagChild[0]] = modelValue.Field(i).Elem().Interface()

				} else {
					mapRes[tagChild[0]] = modelValue.Field(i).Interface()

				}
			}
		}
	}
	return mapRes
}

func ConvertArrayToInterfaceArray[T any](arr []T) []interface{} {
	b := make([]interface{}, len(arr))
	for i := range arr {
		b[i] = arr[i]
	}
	return b
}
