package scenes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/persistence"
	"github.com/paulwrubel/photolum/tracing"
)

// ScenePostHandler handles the /scenes POST endpoint
func ScenePostHandler(response http.ResponseWriter, request *http.Request) {
	// decode request
	var scenePostRequest ScenePostRequest
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
	newSceneID, err := persistence.SaveConfig(newScene)
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
	scenePostResponse := ScenePostResponse{SceneID: newSceneID.String()}
	json.NewEncoder(response).Encode(scenePostResponse)
}

// SceneIDGetHandler handles the /scenes/{scene_id} GET endpoint
func SceneIDGetHandler(response http.ResponseWriter, request *http.Request) {
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
	scene, err := persistence.GetConfig(sceneID)
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
	sceneIDGetResponse := SceneIDGetResponse{Scene: scene}
	json.NewEncoder(response).Encode(sceneIDGetResponse)
}

// SceneIDGetHandler handles the /scenes/{scene_id} GET endpoint
func SceneIDImageGetHandler(response http.ResponseWriter, request *http.Request) {
	// decode request
	var sceneIDImageGetRequest SceneIDImageGetRequest
	err := json.NewDecoder(request.Body).Decode(&sceneIDImageGetRequest)
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
	tracing.SaveImage(sceneID)
	img, err := persistence.GetImage(sceneID)
	if err != nil {
		fmt.Printf("Error retrieving image: %s\n", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Error retrieving image: %s\n", err.Error()),
			},
		})
		return
	}

	switch sceneIDImageGetRequest.Protocol {
	case "image":
		response.Header().Add("Content-Type", "image/png")
		buffer := new(bytes.Buffer)
		err := png.Encode(buffer, img)
		if err != nil {
			fmt.Printf("Error encoding image: %s\n", err.Error())
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(ErrorResponse{
				Error: Error{
					Message: fmt.Sprintf("Error encoding image: %s\n", err.Error()),
				},
			})
		}
		response.WriteHeader(http.StatusOK)
		response.Write(buffer.Bytes())
	case "base64":
		fmt.Printf("Currently Unsupported\n")
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Currently Unsupported\n"),
			},
		})
	default:
		fmt.Printf("Invalid protocol\n")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Invalid Protocol\n"),
			},
		})
		return
	}
}
