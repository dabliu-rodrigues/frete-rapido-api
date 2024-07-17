package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
}

type BadParamError struct {
	Fields    []ParamError `json:"fields"`
	ErrorCode int          `json:"errorCode"`
}

func SendGenericError(w http.ResponseWriter, code int, msg string) {
	response := ErrorResponse{
		ErrorCode: code,
		Message:   msg,
	}
	w.Header().Set("Content-Type", "application/json")

	SendResponse(w, response)
}

func SendBadParamError(w http.ResponseWriter, p []ParamError) {
	response := BadParamError{
		p,
		http.StatusBadRequest,
	}
	SendResponse(w, response)
}

func SendResponse(w http.ResponseWriter, b interface{}) {
	jsonResponse, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
