package imagecontroller

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/enumeration/filetype"
	"github.com/paulwrubel/photolum/persistence/parameterspersistence"
	"github.com/paulwrubel/photolum/persistence/renderpersistence.go"
	"github.com/sirupsen/logrus"
)

var getEndpoint = "/images.GET"

type GetRequest struct {
	RenderName *string `json:"render_name"`
	Format     *string `json:"format"`
}

type GetBase64Response struct {
	Data string `json:"data"`
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
	if getRequest.RenderName == nil ||
		getRequest.Format == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if row exists
	exists, err := renderpersistence.DoesExist(plData, log, *getRequest.RenderName)
	if err != nil {
		errorMessage := "error checking render existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if !exists {
		errorMessage := "render row does not exist"
		errorStatusCode := http.StatusNotFound

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// get render from db
	render, err := renderpersistence.Get(plData, log, *getRequest.RenderName)
	if err != nil {
		errorMessage := "error getting render from database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	switch *getRequest.Format {
	case "base64":
		base64.StdEncoding.EncodeToString(render.ImageData)
		getResponse := GetBase64Response{
			Data: base64.StdEncoding.EncodeToString(render.ImageData),
		}
		response.Header().Add("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(getResponse)
	case "image":
		// get filetype from db
		parameters, err := parameterspersistence.Get(plData, log, render.ParametersName)
		if err != nil {
			errorMessage := "error getting parameters from database"
			errorStatusCode := http.StatusInternalServerError

			log.WithError(err).Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
			return
		}
		switch filetype.FileType(parameters.FileType) {
		case filetype.PNG:
			response.Header().Add("Content-Type", "image/png")
		case filetype.JPEG:
			response.Header().Add("Content-Type", "image/jpeg")
		}
		response.WriteHeader(http.StatusOK)
		response.Write(render.ImageData)
	default:
		errorMessage := "invalid format"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	log.Debug("request completed")
}
