package helloworldcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var endpoint = "/helloworld.GET"

// GetResponse contains the helloworld endpoint response
type GetResponse struct {
	Message string `json:"message"`
}

// GetHandler handles the helloworld GET endpoint
func GetHandler(response http.ResponseWriter, request *http.Request, baseLog *logrus.Logger) {
	requestID, _ := uuid.NewRandom()
	log := baseLog.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"request_id": requestID.String(),
	})
	log.Debug("request received")

	helloWorldResponse := GetResponse{
		Message: "Hello, World!",
	}
	json.NewEncoder(response).Encode(helloWorldResponse)
	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")

	log.Debug("request completed")
}
