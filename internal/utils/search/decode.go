package search

import (
	"encoding/json"
	"errors"
	"go-service/pkg/convert"
	"net/http"
)

func Bind(r *http.Request, filter interface{}) error {
	var mapRes map[string]interface{}
	if r.Method == http.MethodGet {
		mapRes = convert.ConvertObjectToMap(filter)
		queryParams := r.URL.Query()

		for k, v := range queryParams {
			_, exist := mapRes[k]
			if !exist {
				return errors.New("bad request")
			}
			mapRes[k] = v[0]
		}
		jsonBody, err := json.Marshal(mapRes)
		if err != nil {
			return err
		}
		err1 := json.Unmarshal(jsonBody, &filter)

		return err1
	} else {
		err := json.NewDecoder(r.Body).Decode(&filter)
		return err
	}

}
