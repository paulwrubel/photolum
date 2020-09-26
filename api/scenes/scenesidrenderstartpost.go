package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/config/renderstatus"
	"github.com/paulwrubel/photolum/persistence/scene"
	"github.com/paulwrubel/photolum/tracing"
)

// ScenePostHandler handles the /scenes POST endpoint
func ScenesIDRenderStartPostHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData) {
	fmt.Println("Recieved Request for /scenes/{scene_id}/render/start.POST")

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

	doesSceneExist, err := scene.DoesExist(plData, sceneID.String())
	if err != nil {
		fmt.Printf("Error retrieving scene existance status: %s\n", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Error retrieving scene existance status: %s\n", err.Error()),
			},
		})
		return
	}
	if !doesSceneExist {
		fmt.Printf("Error: scene does not exist\n")
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Error: scene does not exist"),
			},
		})
		return
	}

	err = scene.UpdateRenderStatus(plData, sceneID.String(), renderstatus.Pending)
	if err != nil {
		fmt.Printf("Error updating render status: %s\n", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Error updating render status: %s\n", err.Error()),
			},
		})
		return
	}

	tracing.StartRender(plData, sceneID.String())
	if err != nil {
		fmt.Printf("Error starting render: %s\n", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{
			Error: Error{
				Message: fmt.Sprintf("Error starting render: %s\n", err.Error()),
			},
		})
		return
	}

	response.WriteHeader(http.StatusAccepted)
	fmt.Println("Sending Response for /scenes/{scene_id}/render/start.POST")
}
