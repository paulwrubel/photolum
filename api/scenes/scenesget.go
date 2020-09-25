package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paulwrubel/photolum/persistence"
)

// ScenePostResponse contains the scene POST endpoint response
type ScenesGetResponse struct {
	SceneIDs []string `json:"scene_ids"`
}

// ScenePostHandler handles the /scenes POST endpoint
func ScenesGetHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Recieved Request for /scenes.GET")

	sceneDataList := persistence.RetrieveAll()

	sceneIDList := []string{}
	for _, sceneData := range sceneDataList {
		sceneIDList = append(sceneIDList, sceneData.SceneID.String())
	}

	response.WriteHeader(http.StatusCreated)
	scenePostResponse := ScenesGetResponse{SceneIDs: sceneIDList}
	json.NewEncoder(response).Encode(scenePostResponse)
	fmt.Println("Sending Response for /scenes.GET")
}
