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
	mapRes := map[string]interface{}{}
	modelType := reflect.TypeOf(obj)
	modelValue := reflect.ValueOf(obj)
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}
	for i := 0; i < modelType.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldValue := modelValue.Field(i)
		if fieldValue.Kind() == reflect.Struct && fieldType.Anonymous {
			deepToMapOmitEmpty(mapRes, fieldValue.Interface())
		} else {
			tag, ok := fieldType.Tag.Lookup("json")
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

	}
	return mapRes
}

func deepToMapOmitEmpty(structMap map[string]interface{}, obj interface{}) {
	modelType := reflect.TypeOf(obj)
	modelValue := reflect.ValueOf(obj)

	for i := 0; i < modelType.NumField(); i++ {
		fieldValue := modelValue.Field(i)
		fieldType := modelType.Field(i)
		if fieldValue.Kind() == reflect.Struct && fieldType.Anonymous {
			deepToMapOmitEmpty(structMap, fieldValue.Interface())
		} else {
			tag, ok := modelType.Field(i).Tag.Lookup("json")
			if ok {
				tagChild := strings.Split(tag, ",")
				if !modelValue.Field(i).IsZero() {
					if modelValue.Field(i).Kind() == reflect.Pointer {
						structMap[tagChild[0]] = modelValue.Field(i).Elem().Interface()

					} else {
						structMap[tagChild[0]] = modelValue.Field(i).Interface()

					}
				}
			}
		}

	}
}

func ConvertArrayToInterfaceArray[T any](arr []T) []interface{} {
	b := make([]interface{}, len(arr))
	for i := range arr {
		b[i] = arr[i]
	}
	return b
}
