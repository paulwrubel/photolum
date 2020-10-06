package texturecontroller

import (
	"encoding/base64"
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/enumeration/texturetype"
	"github.com/paulwrubel/photolum/persistence/texturepersistence"
	"github.com/paulwrubel/photolum/tracing/shading"
	"github.com/sirupsen/logrus"
)

var getEndpoint = "/textures.GET"
var postEndpoint = "/textures.POST"

type GetRequest struct {
	TextureName *string `json:"texture_name"`
}

type ColorGetResponse struct {
	TextureName string        `json:"texture_name"`
	TextureType string        `json:"texture_type"`
	Color       shading.Color `json:"color"`
}

type ImageGetResponse struct {
	TextureName string  `json:"texture_name"`
	TextureType string  `json:"texture_type"`
	Gamma       float64 `json:"gamma"`
	Magnitude   float64 `json:"magnitude"`
	ImageData   string  `json:"image_data"`
}

type ColorRequest struct {
	Red   *float64 `json:"red"`
	Green *float64 `json:"green"`
	Blue  *float64 `json:"blue"`
}

type PostRequest struct {
	TextureName *string       `json:"texture_name"`
	TextureType *string       `json:"texture_type"`
	Color       *ColorRequest `json:"color"`
	Gamma       *float64      `json:"gamma"`
	Magnitude   *float64      `json:"magnitude"`
	ImageData   *string       `json:"image_data"`
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
	if getRequest.TextureName == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if row exists
	exists, err := texturepersistence.DoesExist(plData, log, *getRequest.TextureName)
	if err != nil {
		errorMessage := "error checking texture existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if !exists {
		errorMessage := "texture row does not exist"
		errorStatusCode := http.StatusNotFound

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// get from db
	texture, err := texturepersistence.Get(plData, log, *getRequest.TextureName)
	if err != nil {
		errorMessage := "error getting texture from database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	var getResponse interface{}
	switch texturetype.TextureType(texture.TextureType) {
	case texturetype.Color:
		getResponse = ColorGetResponse{
			TextureName: texture.TextureName,
			TextureType: texture.TextureType,
			Color: shading.Color{
				Red:   texture.Color[0],
				Green: texture.Color[1],
				Blue:  texture.Color[2],
			},
		}
	case texturetype.Image:
		getResponse = ImageGetResponse{
			TextureName: texture.TextureName,
			TextureType: texture.TextureType,
			Gamma:       *texture.Gamma,
			Magnitude:   *texture.Magnitude,
			ImageData:   base64.StdEncoding.EncodeToString(texture.ImageData),
		}
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
	if postRequest.TextureName == nil ||
		postRequest.TextureType == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// validate input
	errorMessage := ""
	switch texturetype.TextureType(strings.ToUpper(*postRequest.TextureType)) {
	case texturetype.Color:
		if postRequest.Color == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		if *postRequest.Color.Red < 0.0 || *postRequest.Color.Green < 0.0 || *postRequest.Color.Blue < 0.0 {
			errorMessage = "color fields must be greater than or equal to zero"
		}
	case texturetype.Image:
		if postRequest.Gamma == nil ||
			postRequest.Magnitude == nil ||
			postRequest.ImageData == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		if *postRequest.Gamma <= 0.0 {
			errorMessage = "gamma must be greater than zero"
		} else if *postRequest.Magnitude < 0 {
			errorMessage = "magnitude must be greater than or equal to zero"
		} else {
			imageDataReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*postRequest.ImageData))
			_, _, err = image.Decode(imageDataReader)
			if err != nil {
				errMessage := "could not decode image_data"
				errorStatusCode := http.StatusInternalServerError

				log.WithError(err).Error(errorMessage)
				controller.WriteErrorResponse(&response, errorStatusCode, errMessage, err)
				return
			}
		}
	default:
		errorMessage = "invalid texture_type"
	}

	// send error
	if errorMessage != "" {
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if row exists
	exists, err := texturepersistence.DoesExist(plData, log, *postRequest.TextureName)
	if err != nil {
		errorMessage := "error checking texture existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if exists {
		errorMessage := "texture row already exists"
		errorStatusCode := http.StatusConflict

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// assemble texture
	var textureColor []float64
	if postRequest.Color == nil {
		textureColor = nil
	} else {
		textureColor = []float64{*postRequest.Color.Red, *postRequest.Color.Green, *postRequest.Color.Blue}
	}
	var imageData []byte
	if postRequest.ImageData == nil {
		imageData = nil
	} else {
		imageData, _ = base64.StdEncoding.DecodeString(*postRequest.ImageData)
	}
	texture := &texturepersistence.Texture{
		TextureName: *postRequest.TextureName,
		TextureType: strings.ToUpper(*postRequest.TextureType),
		Color:       textureColor,
		Gamma:       postRequest.Gamma,
		Magnitude:   postRequest.Magnitude,
		ImageData:   imageData,
	}

	// save to db
	err = texturepersistence.Save(plData, log, texture)
	if err != nil {
		errorMessage := "error saving texture to database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
	log.Debug("request completed")
}
