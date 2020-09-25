package helloworld

import (
	"encoding/json"
	"net/http"
)

// GetResponse contains the helloworld endpoint response
type GetResponse struct {
	Message string `json:"message"`
}

// GetHandler handles the helloworld GET endpoint
func GetHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	helloWorldResponse := GetResponse{Message: "Hello, World!"}
	json.NewEncoder(response).Encode(helloWorldResponse)
}
