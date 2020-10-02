package parameterscontroller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/enumeration/filetype"
	"github.com/paulwrubel/photolum/persistence/parameterspersistence"
	"github.com/paulwrubel/photolum/tracing/shading"
	"github.com/sirupsen/logrus"
)

var getEndpoint = "/parameters.GET"
var postEndpoint = "/parameters.POST"

type GetRequest struct {
	ParametersName *string `json:"parameters_name"`
}

type BackgroundColorResponse struct {
	Red   float64 `json:"red"`
	Green float64 `json:"green"`
	Blue  float64 `json:"blue"`
}

type GetResponse struct {
	ParametersName           string                  `json:"parameters_name"`
	ImageWidth               int32                   `json:"image_width"`
	ImageHeight              int32                   `json:"image_height"`
	FileType                 filetype.FileType       `json:"file_type"`
	GammaCorrection          float64                 `json:"gamma_correction"`
	TextureGamma             float64                 `json:"texture_gamma"`
	UseScalingTruncation     bool                    `json:"use_scaling_truncation"`
	SamplesPerRound          int32                   `json:"samples_per_round"`
	RoundCount               int32                   `json:"round_count"`
	TileWidth                int32                   `json:"tile_width"`
	TileHeight               int32                   `json:"tile_height"`
	MaxBounces               int32                   `json:"max_bounces"`
	UseBVH                   bool                    `json:"use_bvh"`
	BackgroundColorMagnitude float64                 `json:"background_color_magnitude"`
	BackgroundColor          BackgroundColorResponse `json:"background_color"`
	TMin                     float64                 `json:"t_min"`
	TMax                     float64                 `json:"t_max"`
}

type BackgroundColorRequest struct {
	Red   *float64 `json:"red"`
	Green *float64 `json:"green"`
	Blue  *float64 `json:"blue"`
}

type PostRequest struct {
	ParametersName           *string                 `json:"parameters_name"`
	ImageWidth               *int32                  `json:"image_width"`
	ImageHeight              *int32                  `json:"image_height"`
	FileType                 *filetype.FileType      `json:"file_type"`
	GammaCorrection          *float64                `json:"gamma_correction"`
	TextureGamma             *float64                `json:"texture_gamma"`
	UseScalingTruncation     *bool                   `json:"use_scaling_truncation"`
	SamplesPerRound          *int32                  `json:"samples_per_round"`
	RoundCount               *int32                  `json:"round_count"`
	TileWidth                *int32                  `json:"tile_width"`
	TileHeight               *int32                  `json:"tile_height"`
	MaxBounces               *int32                  `json:"max_bounces"`
	UseBVH                   *bool                   `json:"use_bvh"`
	BackgroundColorMagnitude *float64                `json:"background_color_magnitude"`
	BackgroundColor          *BackgroundColorRequest `json:"background_color"`
	TMin                     *float64                `json:"t_min"`
	TMax                     *float64                `json:"t_max"`
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
	if getRequest.ParametersName == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}
	// check if row exists
	exists, err := parameterspersistence.DoesExist(plData, log, *getRequest.ParametersName)
	if err != nil {
		errorMessage := "error checking parameters existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if !exists {
		errorMessage := "parameters row does not exist"
		errorStatusCode := http.StatusNotFound

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// get from db
	parameters, err := parameterspersistence.Get(plData, log, *getRequest.ParametersName)
	if err != nil {
		errorMessage := "error getting parameters from database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	backgroundColorResponse := BackgroundColorResponse{
		Red:   parameters.BackgroundColor.Red,
		Green: parameters.BackgroundColor.Green,
		Blue:  parameters.BackgroundColor.Blue,
	}
	getResponse := GetResponse{
		ParametersName:           parameters.ParametersName,
		ImageWidth:               parameters.ImageWidth,
		ImageHeight:              parameters.ImageHeight,
		FileType:                 parameters.FileType,
		GammaCorrection:          parameters.GammaCorrection,
		TextureGamma:             parameters.TextureGamma,
		UseScalingTruncation:     parameters.UseScalingTruncation,
		SamplesPerRound:          parameters.SamplesPerRound,
		RoundCount:               parameters.RoundCount,
		TileWidth:                parameters.TileWidth,
		TileHeight:               parameters.TileHeight,
		MaxBounces:               parameters.MaxBounces,
		UseBVH:                   parameters.UseBVH,
		BackgroundColorMagnitude: parameters.BackgroundColorMagnitude,
		BackgroundColor:          backgroundColorResponse,
		TMin:                     parameters.TMin,
		TMax:                     parameters.TMax,
	}
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
	if postRequest.ParametersName == nil ||
		postRequest.ImageWidth == nil ||
		postRequest.ImageHeight == nil ||
		postRequest.FileType == nil ||
		postRequest.GammaCorrection == nil ||
		postRequest.TextureGamma == nil ||
		postRequest.UseScalingTruncation == nil ||
		postRequest.SamplesPerRound == nil ||
		postRequest.RoundCount == nil ||
		postRequest.TileWidth == nil ||
		postRequest.TileHeight == nil ||
		postRequest.MaxBounces == nil ||
		postRequest.UseBVH == nil ||
		postRequest.BackgroundColorMagnitude == nil ||
		postRequest.BackgroundColor == nil ||
		postRequest.TMin == nil ||
		postRequest.TMax == nil ||
		postRequest.BackgroundColor.Red == nil ||
		postRequest.BackgroundColor.Green == nil ||
		postRequest.BackgroundColor.Blue == nil {

		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}
	// check if row exists
	exists, err := parameterspersistence.DoesExist(plData, log, *postRequest.ParametersName)
	if err != nil {
		errorMessage := "error checking parameters existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if exists {
		errorMessage := "parameters row already exists"
		errorStatusCode := http.StatusConflict

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// assemble db row
	backgroundColor := &shading.Color{
		Red:   *(postRequest.BackgroundColor.Red),
		Green: *(postRequest.BackgroundColor.Green),
		Blue:  *(postRequest.BackgroundColor.Blue),
	}
	parameters := &parameterspersistence.Parameters{
		ParametersName:           *(postRequest.ParametersName),
		ImageWidth:               *(postRequest.ImageWidth),
		ImageHeight:              *(postRequest.ImageHeight),
		FileType:                 *(postRequest.FileType),
		GammaCorrection:          *(postRequest.GammaCorrection),
		TextureGamma:             *(postRequest.TextureGamma),
		UseScalingTruncation:     *(postRequest.UseScalingTruncation),
		SamplesPerRound:          *(postRequest.SamplesPerRound),
		RoundCount:               *(postRequest.RoundCount),
		TileWidth:                *(postRequest.TileWidth),
		TileHeight:               *(postRequest.TileHeight),
		MaxBounces:               *(postRequest.MaxBounces),
		UseBVH:                   *(postRequest.UseBVH),
		BackgroundColorMagnitude: *(postRequest.BackgroundColorMagnitude),
		BackgroundColor:          backgroundColor,
		TMin:                     *(postRequest.TMin),
		TMax:                     *(postRequest.TMax),
	}

	// save to db
	err = parameterspersistence.Save(plData, log, parameters)
	if err != nil {
		errorMessage := "error saving parameters to database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
	log.Debug("request completed")
}
