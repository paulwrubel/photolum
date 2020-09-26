package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/persistence/scene"
)

// ScenePostResponse contains the scene POST endpoint response
type ScenesGetResponse struct {
	SceneIDs []string `json:"scene_ids"`
}

// ScenePostHandler handles the /scenes POST endpoint
func ScenesGetHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData) {
	fmt.Println("Recieved Request for /scenes.GET")

	sceneList, err := scene.RetrieveAll(plData)
	if err != nil {
		fmt.Printf("Error retrieving all scenes: %s\n", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Error retrieving all scenes: %s\n", err.Error()),
			},
		})
	}

	sceneIDList := []string{}
	for _, scn := range sceneList {
		sceneIDList = append(sceneIDList, scn.SceneID)
	}

	response.WriteHeader(http.StatusCreated)
	scenePostResponse := ScenesGetResponse{SceneIDs: sceneIDList}
	json.NewEncoder(response).Encode(scenePostResponse)
	fmt.Println("Sending Response for /scenes.GET")
}
