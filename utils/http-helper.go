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

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func SendBadParamError(w http.ResponseWriter, p []ParamError) {
	response := BadParamError{
		p,
		http.StatusBadRequest,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonResponse)
}
