package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/persistence/scene"
)

// SceneIDGetRequest contains the sceneID GET endpoint request
type SceneIDGetResponse struct {
	Scene scene.Scene `json:"scene"`
}

// SceneIDGetHandler handles the /scenes/{scene_id} GET endpoint
func SceneIDGetHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData) {
	fmt.Println("Recieved Request for /scenes/{scene_id}.GET")
	// decode request
	params := mux.Vars(request)
	// assemble and save scene
	sceneID, err := uuid.Parse(params["scene_id"])
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
	retrievedScene, err := scene.Retrieve(plData, sceneID.String())
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

	response.WriteHeader(http.StatusOK)
	sceneIDGetResponse := SceneIDGetResponse{Scene: *retrievedScene}
	json.NewEncoder(response).Encode(sceneIDGetResponse)
	fmt.Println("Sending Response for /scenes/{scene_id}.GET")
}
