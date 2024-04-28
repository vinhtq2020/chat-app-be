package response

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, statusCode int, body interface{}) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}
