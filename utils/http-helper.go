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

	SendResponse(w, code, response)
}

func SendBadParamError(w http.ResponseWriter, p []ParamError) {
	response := BadParamError{
		p,
		http.StatusBadRequest,
	}
	SendResponse(w, http.StatusBadRequest, response)
}

func SendOKResponse(w http.ResponseWriter, b interface{}) {
	SendResponse(w, http.StatusOK, b)
}

func SendResponse(w http.ResponseWriter, code int, b interface{}) {
	jsonResponse, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}
