package cameracontroller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/constants"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/persistence/camerapersistence"
	"github.com/paulwrubel/photolum/tracing/geometry"
	"github.com/sirupsen/logrus"
)

var getEndpoint = "/cameras.GET"
var postEndpoint = "/cameras.POST"

type GetRequest struct {
	CameraName *string `json:"camera_name"`
}

type GetResponse struct {
	CameraName     string          `json:"camera_name"`
	EyeLocation    geometry.Point  `json:"eye_location"`
	TargetLocation geometry.Point  `json:"target_location"`
	UpVector       geometry.Vector `json:"up_vector"`
	VerticalFOV    float64         `json:"vertical_fov"`
	Aperture       float64         `json:"aperture"`
	FocusDistance  float64         `json:"focus_distance"`
}

type VectorRequest struct {
	X *float64 `json:"x"`
	Y *float64 `json:"y"`
	Z *float64 `json:"z"`
}

type PostRequest struct {
	CameraName     *string        `json:"camera_name"`
	EyeLocation    *VectorRequest `json:"eye_location"`
	TargetLocation *VectorRequest `json:"target_location"`
	UpVector       *VectorRequest `json:"up_vector"`
	VerticalFOV    *float64       `json:"vertical_fov"`
	Aperture       *float64       `json:"aperture"`
	FocusDistance  *float64       `json:"focus_distance"`
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
	if getRequest.CameraName == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}
	// check if row exists
	exists, err := camerapersistence.DoesExist(plData, log, *getRequest.CameraName)
	if err != nil {
		errorMessage := "error checking camera existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if !exists {
		errorMessage := "camera row does not exist"
		errorStatusCode := http.StatusNotFound

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// get from db
	camera, err := camerapersistence.Get(plData, log, *getRequest.CameraName)
	if err != nil {
		errorMessage := "error getting camera from database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	getResponse := GetResponse{
		CameraName:     camera.CameraName,
		EyeLocation:    camera.EyeLocation,
		TargetLocation: camera.TargetLocation,
		UpVector:       camera.UpVector,
		VerticalFOV:    camera.VerticalFOV,
		Aperture:       camera.Aperture,
		FocusDistance:  camera.FocusDistance,
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
	if postRequest.CameraName == nil ||
		postRequest.EyeLocation == nil ||
		postRequest.EyeLocation.X == nil ||
		postRequest.EyeLocation.Y == nil ||
		postRequest.EyeLocation.Z == nil ||
		postRequest.TargetLocation == nil ||
		postRequest.TargetLocation.X == nil ||
		postRequest.TargetLocation.Y == nil ||
		postRequest.TargetLocation.Z == nil ||
		postRequest.UpVector == nil ||
		postRequest.UpVector.X == nil ||
		postRequest.UpVector.Y == nil ||
		postRequest.UpVector.Z == nil ||
		postRequest.VerticalFOV == nil ||
		postRequest.Aperture == nil ||
		postRequest.FocusDistance == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// validate input
	var errorMessage = ""
	if *postRequest.VerticalFOV < constants.CameraMinimumVerticalFOV {
		errorMessage = fmt.Sprintf("vertical_fov cannot be below %d", constants.CameraMinimumVerticalFOV)
	}
	if *postRequest.VerticalFOV > constants.CameraMaximumVerticalFOV {
		errorMessage = fmt.Sprintf("vertical_fov cannot exceed %d", constants.CameraMaximumVerticalFOV)
	}
	if *postRequest.Aperture < constants.CameraMinimumAperture {
		errorMessage = fmt.Sprintf("aperture cannot be below %d", constants.CameraMinimumAperture)
	}
	if *postRequest.Aperture > constants.CameraMaximumAperture {
		errorMessage = fmt.Sprintf("aperture cannot exceed %d", constants.CameraMaximumAperture)
	}
	if *postRequest.FocusDistance < constants.CameraMinimumFocusDistance {
		errorMessage = fmt.Sprintf("focus_distance cannot be below %d", constants.CameraMinimumFocusDistance)
	}
	if *postRequest.FocusDistance > constants.CameraMaximumFocusDistance {
		errorMessage = fmt.Sprintf("focus_distance cannot exceed %d", constants.CameraMaximumFocusDistance)
	}

	// send error
	if errorMessage != "" {
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if row exists
	exists, err := camerapersistence.DoesExist(plData, log, *postRequest.CameraName)
	if err != nil {
		errorMessage := "error checking camera existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if exists {
		errorMessage := "camera row already exists"
		errorStatusCode := http.StatusConflict

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// assemble camera
	camera := &camerapersistence.Camera{
		CameraName: *postRequest.CameraName,
		EyeLocation: geometry.Point{
			X: *postRequest.EyeLocation.X,
			Y: *postRequest.EyeLocation.Y,
			Z: *postRequest.EyeLocation.Z,
		},
		TargetLocation: geometry.Point{
			X: *postRequest.TargetLocation.X,
			Y: *postRequest.TargetLocation.Y,
			Z: *postRequest.TargetLocation.Z,
		},
		UpVector: geometry.Vector{
			X: *postRequest.UpVector.X,
			Y: *postRequest.UpVector.Y,
			Z: *postRequest.UpVector.Z,
		},
		VerticalFOV:   *postRequest.VerticalFOV,
		Aperture:      *postRequest.Aperture,
		FocusDistance: *postRequest.FocusDistance,
	}

	// save to db
	err = camerapersistence.Save(plData, log, camera)
	if err != nil {
		errorMessage := "error saving camera to database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
	log.Debug("request completed")
}
