package response

import (
	"encoding/json"
	"net/http"
)

type (
	responseError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func newResponseError(err error, statusCode int) responseError {
	return responseError{Code: statusCode, Message: err.Error()}
}

func StatusCode(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func JSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(&data)
}

func Error(w http.ResponseWriter, err error, statusCode int) {
	JSON(w, newResponseError(statusCode, err), statusCode)
}
