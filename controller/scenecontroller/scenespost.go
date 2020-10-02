package scenecontroller

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"

// 	"github.com/paulwrubel/photolum/config"
// 	"github.com/paulwrubel/photolum/controller"
// 	"github.com/paulwrubel/photolum/persistence/scenepersistence"
// 	"github.com/paulwrubel/photolum/service/sceneservice"
// )

// // ScenePostRequest contains the scene POST endpoint request
// type ScenesPostRequest struct {
// 	ImageWidth     int      `json:"image_width"`      // width of the image in pixels
// 	ImageHeight    int      `json:"image_height"`     // height of the image in pixels
// 	ImageFileTypes []string `json:"image_file_types"` // image file type (png, jpg, etc.)
// }

// // ScenePostResponse contains the scene POST endpoint response
// type ScenesPostResponse struct {
// 	SceneID string `json:"scene_id"`
// }

// // ScenePostHandler handles the /scenes POST endpoint
// func ScenesPostHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData) {
// 	fmt.Println("Recieved Request for /scenes.POST")

// 	// decode request
// 	var scenePostRequest ScenesPostRequest
// 	if request.Body != nil {
// 		defer request.Body.Close()
// 	}
// 	err := json.NewDecoder(request.Body).Decode(&scenePostRequest)
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("Error decoding request: %s", err.Error())
// 		errorStatusCode := http.StatusBadRequest
// 		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
// 		return
// 	}

// 	// get file types and validate/remove duplicates
// 	fileTypeMap := make(map[string]int)
// 	for _, fileType := range scenePostRequest.ImageFileTypes {
// 		if fileType == "png" {
// 			fileTypeMap["png"] = 0
// 		} else if fileType == "jpeg" || fileType == "jpg" {
// 			fileTypeMap["jpeg"] = 0
// 		} else {
// 			errorMessage := fmt.Sprint("Invalid image file type")
// 			errorStatusCode := http.StatusBadRequest
// 			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
// 			return
// 		}
// 	}

// 	var formattedFileTypes []string
// 	for key := range fileTypeMap {
// 		formattedFileTypes = append(formattedFileTypes, key)
// 	}

// 	// assemble and save scene
// 	scn := &scenepersistence.Scene{
// 		ImageWidth:     scenePostRequest.ImageWidth,
// 		ImageHeight:    scenePostRequest.ImageHeight,
// 		ImageFileTypes: strings.Join(formattedFileTypes, ","),
// 	}
// 	sceneID, err := sceneservice.Save(plData, scn)
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("Error saving scene: %s", err.Error())
// 		errorStatusCode := http.StatusInternalServerError
// 		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
// 		return
// 	}

// 	response.WriteHeader(http.StatusCreated)
// 	response.Header().Add("Content-Type", "application/json")
// 	json.NewEncoder(response).Encode(ScenesPostResponse{
// 		SceneID: sceneID,
// 	})
// 	fmt.Println("Sending Response for /scenes.POST")
// }
