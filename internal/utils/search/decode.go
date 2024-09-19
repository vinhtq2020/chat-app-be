package search

import (
	"encoding/json"
	"net/http"
)

func Bind(r *http.Request, filter interface{}) error {
	mapRes := map[string]interface{}{}
	if r.Method == http.MethodGet {
		queryParams := r.URL.Query()

		for k, v := range queryParams {
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
