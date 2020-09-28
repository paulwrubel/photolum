package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorInfo struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	ErrorInfo ErrorInfo `json:"error"`
}

func WriteErrorResponse(response *http.ResponseWriter, statusCode int, errorMessage string) {
	fmt.Println(errorMessage)
	(*response).Header().Add("Content-Type", "application/json")
	(*response).WriteHeader(statusCode)
	json.NewEncoder(*response).Encode(ErrorResponse{
		ErrorInfo: ErrorInfo{
			Message: fmt.Sprint(errorMessage),
		},
	})
}
