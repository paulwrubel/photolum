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
	_image "github.com/paulwrubel/photolum/persistence/image"
)

// SceneIDGetRequest contains the sceneID GET endpoint request
type SceneIDImageGetRequest struct {
	Protocol string `json:"protocol"`
}

// SceneIDGetBase64Response contains the sceneID GET endpoint request in base64 format
type SceneIDImageGetBase64Response struct {
	Scene string `json:"scene"`
}

// SceneIDGetHandler handles the /scenes/{scene_id} GET endpoint
func SceneIDImageGetHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData) {
	fmt.Println("Recieved Request for /scenes/{scene_id}/image.GET")
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
	image, err := _image.Retrieve(plData, sceneID.String())
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
	img := image.ImageData
	if img == nil {
		response.WriteHeader(http.StatusNotFound)
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
	fmt.Println("Sending Response for /scenes/{scene_id}/image.GET")
}
