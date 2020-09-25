package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/persistence"
)

// ScenePostRequest contains the scene POST endpoint request
type ScenesPostRequest struct {
	Scene config.Scene `json:"scene"`
}

// ScenePostResponse contains the scene POST endpoint response
type ScenesPostResponse struct {
	SceneID string `json:"scene_id"`
}

// ScenePostHandler handles the /scenes POST endpoint
func ScenesPostHandler(response http.ResponseWriter, request *http.Request) {
	// decode request
	var scenePostRequest ScenesPostRequest
	err := json.NewDecoder(request.Body).Decode(&scenePostRequest)
	if err != nil {
		fmt.Printf("Error decoding request: %s\n", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Error decoding request: %s\n", err.Error()),
			},
		})
		return
	}
	defer request.Body.Close()

	// assemble and save scene
	newScene := config.Scene{
		ImageWidth:  scenePostRequest.Scene.ImageWidth,
		ImageHeight: scenePostRequest.Scene.ImageHeight,
		FileType:    scenePostRequest.Scene.FileType,
	}
	sceneData := persistence.SceneData{
		Scene: newScene,
	}
	newSceneID, err := persistence.Create(sceneData)
	if err != nil {
		fmt.Printf("Error saving scene: %s\n", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Error saving scene: %s\n", err.Error()),
			},
		})
		return
	}

	response.WriteHeader(http.StatusCreated)
	scenePostResponse := ScenesPostResponse{SceneID: newSceneID.String()}
	json.NewEncoder(response).Encode(scenePostResponse)
}
