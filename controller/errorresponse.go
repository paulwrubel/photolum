package controller

import (
	"encoding/json"
	"net/http"
)

type ErrorInfo struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type ErrorResponse struct {
	ErrorInfo ErrorInfo `json:"error"`
}

func WriteErrorResponse(response *http.ResponseWriter, statusCode int, errorMessage string, err error) {
	(*response).Header().Add("Content-Type", "application/json")
	(*response).WriteHeader(statusCode)
	errorString := ""
	if err != nil {
		errorString = err.Error()
	}
	json.NewEncoder(*response).Encode(ErrorResponse{
		ErrorInfo: ErrorInfo{
			Message: errorMessage,
			Error:   errorString,
		},
	})
}
