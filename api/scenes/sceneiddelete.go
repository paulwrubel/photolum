package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paulwrubel/photolum/persistence"
)

// ScenePostHandler handles the /scenes POST endpoint
func SceneIDDeleteHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Recieved Request for /scenes/{scene_id}.DELETE")

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

	persistence.Delete(sceneID)

	response.WriteHeader(http.StatusNoContent)
	fmt.Println("Sending Response for /scenes/{scene_id}.DELETE")
}
