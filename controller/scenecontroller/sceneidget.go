package scenecontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/service/sceneservice"
)

// SceneIDGetRequest contains the sceneID GET endpoint request
type SceneIDGetResponse struct {
	ImageWidth     int      `json:"image_width"`      // width of the image in pixels
	ImageHeight    int      `json:"image_height"`     // height of the image in pixels
	ImageFileTypes []string `json:"image_file_types"` // image file type (png, jpg, etc.)
}

// SceneIDGetHandler handles the /scenes/{scene_id} GET endpoint
func SceneIDGetHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData) {
	fmt.Println("Recieved Request for /scenes/{scene_id}.GET")

	// getting request info
	params := mux.Vars(request)
	sceneID := params["scene_id"]

	// validating uuid
	_, err := uuid.Parse(sceneID)
	if err != nil {
		errorMessage := fmt.Sprintf("UUID is not valid (malformatted): %s", err.Error())
		errorStatusCode := http.StatusBadRequest
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}

	// checking existence
	sceneDoesExist, err := sceneservice.DoesExist(plData, sceneID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error checking Scene existance: %s", err.Error())
		errorStatusCode := http.StatusInternalServerError
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}
	if !sceneDoesExist {
		errorMessage := fmt.Sprint("Error: Scene does not exist")
		errorStatusCode := http.StatusNotFound
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}

	// get scene
	scn, err := sceneservice.Get(plData, sceneID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error retrieving Scene: %s", err.Error())
		errorStatusCode := http.StatusInternalServerError
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	sceneIDGetResponse := SceneIDGetResponse{
		ImageWidth:     scn.ImageWidth,
		ImageHeight:    scn.ImageHeight,
		ImageFileTypes: strings.Split(scn.ImageFileTypes, ","),
	}
	json.NewEncoder(response).Encode(sceneIDGetResponse)
	fmt.Println("Sending Response for /scenes/{scene_id}.GET")
}
