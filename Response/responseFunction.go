package Response

import (
	"encoding/json"
	"net/http"
)

func ShowResponse(status string, statusCode int64,message string, data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(int(statusCode))
	response := Response{
		Status: status,
		Code:   statusCode,
		Message: message,
		Data:   data,
	}

	json.NewEncoder(w).Encode(&response)
}
