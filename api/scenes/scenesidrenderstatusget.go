package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/persistence"
)

// ScenePostResponse contains the scene POST endpoint response
type SceneIDRenderStatusGetResponse struct {
	RenderStatus string `json:"render_status"`
}

// ScenePostHandler handles the /scenes POST endpoint
func ScenesIDRenderStatusGetHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Recieved Request for /scenes/{scene_id}/render/status.GET")

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

	sceneData, err := persistence.Retrieve(sceneID)
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

	response.WriteHeader(http.StatusCreated)
	sceneIDRenderStatusGetResponse := SceneIDRenderStatusGetResponse{RenderStatus: string(sceneData.RenderStatus)}
	json.NewEncoder(response).Encode(sceneIDRenderStatusGetResponse)
	fmt.Println("Sending Response for /scenes/{scene_id}/render/status.GET")
}