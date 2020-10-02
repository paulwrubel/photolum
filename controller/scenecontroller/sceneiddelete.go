package scenecontroller

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/google/uuid"
// 	"github.com/gorilla/mux"
// 	"github.com/paulwrubel/photolum/config"
// 	"github.com/paulwrubel/photolum/controller"
// 	"github.com/paulwrubel/photolum/service/sceneservice"
// )

// // ScenePostHandler handles the /scenes POST endpoint
// func SceneIDDeleteHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData) {
// 	fmt.Println("Recieved Request for /scenes/{scene_id}.DELETE")

// 	// getting request info
// 	params := mux.Vars(request)
// 	sceneID := params["scene_id"]

// 	// validating uuid
// 	_, err := uuid.Parse(sceneID)
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("UUID is not valid (malformatted): %s", err.Error())
// 		errorStatusCode := http.StatusBadRequest
// 		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
// 		return
// 	}

// 	// checking existence
// 	sceneDoesExist, err := sceneservice.DoesExist(plData, sceneID)
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("Error checking Scene existance: %s", err.Error())
// 		errorStatusCode := http.StatusInternalServerError
// 		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
// 		return
// 	}
// 	if !sceneDoesExist {
// 		errorMessage := fmt.Sprint("Error: Scene does not exist")
// 		errorStatusCode := http.StatusNotFound
// 		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
// 		return
// 	}

// 	// do the delete
// 	err = sceneservice.Delete(plData, sceneID)
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("Error deleting scene: %s", err.Error())
// 		errorStatusCode := http.StatusInternalServerError
// 		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage)
// 		return
// 	}

// 	response.WriteHeader(http.StatusNoContent)
// 	fmt.Println("Sending Response for /scenes/{scene_id}.DELETE")
// }
