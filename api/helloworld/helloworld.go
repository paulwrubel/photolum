package helloworld

import (
	"encoding/json"
	"net/http"
)

// GetHandler handles the helloworld GET endpoint
func GetHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	helloWorldResponse := GetResponse{Message: "Hello, World!"}
	json.NewEncoder(response).Encode(helloWorldResponse)
}
