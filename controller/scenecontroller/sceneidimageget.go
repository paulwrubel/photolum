package scenecontroller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/service/imageservice"
	"github.com/paulwrubel/photolum/service/sceneservice"
)

// SceneIDGetRequest contains the sceneID GET endpoint request
type SceneIDImageGetRequest struct {
	Protocol string `json:"protocol"`
	FileType string `json:"file_type"`
}

// SceneIDGetBase64Response contains the sceneID GET endpoint request in base64 format
type SceneIDImageGetBase64Response struct {
	ImageData string `json:"image_data"`
}

// SceneIDGetHandler handles the /scenes/{scene_id} GET endpoint
func SceneIDImageGetHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData) {
	fmt.Println("Recieved Request for /scenes/{scene_id}/image.GET")

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

	// decode request
	var sceneIDImageGetRequest SceneIDImageGetRequest
	if request.Body != nil {
		defer request.Body.Close()
	}
	err = json.NewDecoder(request.Body).Decode(&sceneIDImageGetRequest)
	if err != nil {
		errorMessage := fmt.Sprintf("Error decoding request: %s", err.Error())
		errorStatusCode := http.StatusBadRequest
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}

	var formattedFileType string
	if sceneIDImageGetRequest.FileType == "png" {
		formattedFileType = "png"
	} else if sceneIDImageGetRequest.FileType == "jpeg" || sceneIDImageGetRequest.FileType == "jpg" {
		formattedFileType = "jpeg"
	} else {
		errorMessage := fmt.Sprint("Invalid image file type")
		errorStatusCode := http.StatusBadRequest
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}

	// checking image existence
	imageDoesExist, err := imageservice.DoesExist(plData, sceneID, formattedFileType)
	if err != nil {
		errorMessage := fmt.Sprintf("Error checking Image existance: %s", err.Error())
		errorStatusCode := http.StatusInternalServerError
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}
	if !imageDoesExist {
		errorMessage := fmt.Sprint("Error: Image does not exist")
		errorStatusCode := http.StatusNotFound
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}

	// get image
	imgBytes, err := imageservice.GetEncoded(plData, sceneID, formattedFileType)
	if err != nil {
		errorMessage := fmt.Sprintf("Error retrieving image: %s", err.Error())
		errorStatusCode := http.StatusInternalServerError
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}

	switch sceneIDImageGetRequest.Protocol {
	case "image":
		response.Header().Add("Content-Type", "image/"+formattedFileType)
		response.WriteHeader(http.StatusOK)
		response.Write(imgBytes)
	case "base64":
		base64EncodedString := base64.StdEncoding.EncodeToString(imgBytes)
		response.Header().Add("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(SceneIDImageGetBase64Response{
			ImageData: base64EncodedString,
		})
	default:
		errorMessage := fmt.Sprint("Invalid protocol")
		errorStatusCode := http.StatusBadRequest
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
		return
	}
	fmt.Println("Sending Response for /scenes/{scene_id}/image.GET")
}
