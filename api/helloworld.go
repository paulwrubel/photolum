package api

import (
	"encoding/json"
	"net/http"
)

// HelloWorldResponse contains the helloworld endpoint response
type HelloWorldResponse struct {
	Message string `json:"message"`
}

// HelloWorldHandler handles the helloworld endpoint
func HelloWorldHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	helloWorldResponse := HelloWorldResponse{Message: "Hello, World!"}
	json.NewEncoder(response).Encode(helloWorldResponse)
}
