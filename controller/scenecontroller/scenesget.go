package scenecontroller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/service/sceneservice"
)

// ScenePostResponse contains the scene POST endpoint response
type ScenesGetResponse struct {
	SceneIDs []string `json:"scene_ids"`
}

// ScenePostHandler handles the /scenes POST endpoint
func ScenesGetHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData) {
	fmt.Println("Recieved Request for /scenes.GET")

	sceneList, err := sceneservice.GetAll(plData)
	if err != nil {
		errorMessage := fmt.Sprintf("Error retrieving all scenes: %s", err.Error())
		errorStatusCode := http.StatusBadRequest
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}

	sceneIDList := []string{}
	for _, scn := range sceneList {
		sceneIDList = append(sceneIDList, scn.SceneID)
	}

	response.WriteHeader(http.StatusCreated)
	response.Header().Add("Content-Type", "application/json")
	scenePostResponse := ScenesGetResponse{SceneIDs: sceneIDList}
	json.NewEncoder(response).Encode(scenePostResponse)
	fmt.Println("Sending Response for /scenes.GET")
}
