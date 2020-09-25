package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/persistence"
	"github.com/paulwrubel/photolum/tracing"
)

// ScenePostHandler handles the /scenes POST endpoint
func ScenesIDRenderStopPostHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Recieved Request for /scenes/{scene_id}/render/stop.POST")

	params := mux.Vars(request)
	sceneIDString := params["scene_id"]
	sceneID, err := uuid.Parse(sceneIDString)
	if err != nil {
		fmt.Printf("Error parsing uuid: %s\n", err.Error())
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Error parsing uuid: %s\n", err.Error()),
			},
		})
		return
	}

	_, err = persistence.Retrieve(sceneID)
	if err != nil {
		fmt.Printf("Error retrieving scene: %s\n", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Error retrieving scene: %s\n", err.Error()),
			},
		})
		return
	}

	tracing.StopRender(sceneID)

	response.WriteHeader(http.StatusAccepted)
	fmt.Println("Sending Response for /scenes/{scene_id}/render/stop.POST")
}