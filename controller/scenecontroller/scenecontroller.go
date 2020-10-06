package scenecontroller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/persistence/camerapersistence"
	"github.com/paulwrubel/photolum/persistence/materialpersistence"
	"github.com/paulwrubel/photolum/persistence/primitivepersistence"
	"github.com/paulwrubel/photolum/persistence/scenepersistence"
	"github.com/paulwrubel/photolum/persistence/sceneprimitivematerialpersistence"
	"github.com/sirupsen/logrus"
)

var getEndpoint = "/scenes.GET"
var postEndpoint = "/scenes.POST"

type GetRequest struct {
	SceneName *string `json:"scene_name"`
}

type PrimitiveMaterialGetResponse struct {
	PrimitiveName string `json:"primitive_name"`
	MaterialName  string `json:"material_name"`
}

type GetResponse struct {
	SceneName          string                         `json:"scene_name"`
	CameraName         string                         `json:"camera_name"`
	PrimitiveMaterials []PrimitiveMaterialGetResponse `json:"primitive_materials"`
}

type PrimitiveMaterialPostRequest struct {
	PrimitiveName *string `json:"primitive_name"`
	MaterialName  *string `json:"material_name"`
}

type PostRequest struct {
	SceneName          *string                        `json:"scene_name"`
	CameraName         *string                        `json:"camera_name"`
	PrimitiveMaterials []PrimitiveMaterialGetResponse `json:"primitive_materials"`
}

func GetHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData, baseLog *logrus.Logger) {
	requestID, _ := uuid.NewRandom()
	log := baseLog.WithFields(logrus.Fields{
		"endpoint":   getEndpoint,
		"request_id": requestID.String(),
	})
	log.Debug("request received")

	// decode request
	var getRequest *GetRequest
	if request.Body != nil {
		defer request.Body.Close()
	}
	err := json.NewDecoder(request.Body).Decode(&getRequest)
	if err != nil {
		errorMessage := "error decoding request body"
		errorStatusCode := http.StatusBadRequest

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	// check for missing fields
	if getRequest.SceneName == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if row exists
	exists, err := scenepersistence.DoesExist(plData, log, *getRequest.SceneName)
	if err != nil {
		errorMessage := "error checking scene existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if !exists {
		errorMessage := "scene row does not exist"
		errorStatusCode := http.StatusNotFound

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// get top-level scene from db
	scene, err := scenepersistence.Get(plData, log, *getRequest.SceneName)
	if err != nil {
		errorMessage := "error getting scene from database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	// get attached sceneprimitivematerials from db
	scenePrimitiveMaterials, err := sceneprimitivematerialpersistence.GetAllInScene(plData, log, *getRequest.SceneName)
	if err != nil {
		errorMessage := "error getting scene from database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	primitiveMaterialGetResponses := []PrimitiveMaterialGetResponse{}
	for _, spm := range scenePrimitiveMaterials {
		pmgr := PrimitiveMaterialGetResponse{
			PrimitiveName: spm.PrimitiveName,
			MaterialName:  spm.MaterialName,
		}
		primitiveMaterialGetResponses = append(primitiveMaterialGetResponses, pmgr)
	}

	getResponse := GetResponse{
		SceneName:          scene.SceneName,
		CameraName:         scene.CameraName,
		PrimitiveMaterials: primitiveMaterialGetResponses,
	}
	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(getResponse)

	log.Debug("request completed")
}

func PostHandler(response http.ResponseWriter, request *http.Request, plData *config.PhotolumData, baseLog *logrus.Logger) {
	requestID, _ := uuid.NewRandom()
	log := baseLog.WithFields(logrus.Fields{
		"endpoint":   postEndpoint,
		"request_id": requestID.String(),
	})
	log.Debug("request received")

	// decode request
	var postRequest *PostRequest
	if request.Body != nil {
		defer request.Body.Close()
	}
	err := json.NewDecoder(request.Body).Decode(&postRequest)
	if err != nil {
		errorMessage := "error decoding request body"
		errorStatusCode := http.StatusBadRequest

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	// check for missing fields
	if postRequest.SceneName == nil ||
		postRequest.CameraName == nil ||
		postRequest.PrimitiveMaterials == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// validate input
	var errorMessage = ""
	if len(postRequest.PrimitiveMaterials) == 0 {
		errorMessage = "scene must contain at least one primitive_material"
	}

	// check if camera exists
	exists, err := camerapersistence.DoesExist(plData, log, *postRequest.CameraName)
	if err != nil {
		errorMessage := "error checking camera existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if !exists {
		errorMessage = "named camera does not exist"
	}

	// check if primitives and materials exists
	for index, spm := range postRequest.PrimitiveMaterials {
		// check primitive
		exists, err := primitivepersistence.DoesExist(plData, log, spm.PrimitiveName)
		if err != nil {
			errorMessage := "error checking primitive existence in database"
			errorStatusCode := http.StatusInternalServerError

			log.WithError(err).Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
			return
		}
		if !exists {
			errorMessage = fmt.Sprintf("named primitive at index: %d does not exist", index)
			break
		}
		// check material
		exists, err = materialpersistence.DoesExist(plData, log, spm.MaterialName)
		if err != nil {
			errorMessage := "error checking material existence in database"
			errorStatusCode := http.StatusInternalServerError

			log.WithError(err).Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
			return
		}
		if !exists {
			errorMessage = fmt.Sprintf("named material at index: %d does not exist", index)
			break
		}
	}

	// send error
	if errorMessage != "" {
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if scene row exists
	exists, err := scenepersistence.DoesExist(plData, log, *postRequest.SceneName)
	if err != nil {
		errorMessage := "error checking scene existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if exists {
		errorMessage := "scene row already exists"
		errorStatusCode := http.StatusConflict

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if sceneprimitivematerial rows exist
	for _, spm := range postRequest.PrimitiveMaterials {
		exists, err := sceneprimitivematerialpersistence.DoesExist(plData, log, *postRequest.SceneName, spm.PrimitiveName, spm.PrimitiveName)
		if err != nil {
			errorMessage := "error checking sceneprimitivematerial existence in database"
			errorStatusCode := http.StatusInternalServerError

			log.WithError(err).Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
			return
		}
		if exists {
			errorMessage := "sceneprimitivematerial row already exists"
			errorStatusCode := http.StatusConflict

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
	}

	// assemble scene
	scene := &scenepersistence.Scene{
		SceneName:  *postRequest.SceneName,
		CameraName: *postRequest.CameraName,
	}
	// assemble sceneprimitivematerials
	scenePrimitiveMaterials := []*sceneprimitivematerialpersistence.ScenePrimitiveMaterial{}
	for _, spm := range postRequest.PrimitiveMaterials {
		scenePrimitiveMaterial := &sceneprimitivematerialpersistence.ScenePrimitiveMaterial{
			SceneName:     *postRequest.SceneName,
			PrimitiveName: spm.PrimitiveName,
			MaterialName:  spm.MaterialName,
		}
		scenePrimitiveMaterials = append(scenePrimitiveMaterials, scenePrimitiveMaterial)
	}

	// save scene to db
	err = scenepersistence.Save(plData, log, scene)
	if err != nil {
		errorMessage := "error saving scene to database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	// save sceneprimitivematerials to db
	for _, spm := range scenePrimitiveMaterials {
		err = sceneprimitivematerialpersistence.Save(plData, log, spm)
		if err != nil {
			errorMessage := "error saving sceneprimitivematerial to database"
			errorStatusCode := http.StatusInternalServerError

			log.WithError(err).Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
			return
		}
	}

	response.WriteHeader(http.StatusCreated)
	log.Debug("request completed")
}
