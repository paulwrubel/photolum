package parameterscontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/constants"
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

type GetResponse struct {
	ParametersName           string        `json:"parameters_name"`
	ImageWidth               uint32        `json:"image_width"`
	ImageHeight              uint32        `json:"image_height"`
	FileType                 string        `json:"file_type"`
	GammaCorrection          float64       `json:"gamma_correction"`
	TextureGamma             float64       `json:"texture_gamma"`
	UseScalingTruncation     bool          `json:"use_scaling_truncation"`
	SamplesPerRound          uint32        `json:"samples_per_round"`
	RoundCount               uint32        `json:"round_count"`
	TileWidth                uint32        `json:"tile_width"`
	TileHeight               uint32        `json:"tile_height"`
	MaxBounces               uint32        `json:"max_bounces"`
	UseBVH                   bool          `json:"use_bvh"`
	BackgroundColorMagnitude float64       `json:"background_color_magnitude"`
	BackgroundColor          shading.Color `json:"background_color"`
	TMin                     float64       `json:"t_min"`
	TMax                     float64       `json:"t_max"`
}

type ColorRequest struct {
	Red   *float64 `json:"red"`
	Green *float64 `json:"green"`
	Blue  *float64 `json:"blue"`
}

type PostRequest struct {
	ParametersName           *string       `json:"parameters_name"`
	ImageWidth               *uint32       `json:"image_width"`
	ImageHeight              *uint32       `json:"image_height"`
	FileType                 *string       `json:"file_type"`
	GammaCorrection          *float64      `json:"gamma_correction"`
	TextureGamma             *float64      `json:"texture_gamma"`
	UseScalingTruncation     *bool         `json:"use_scaling_truncation"`
	SamplesPerRound          *uint32       `json:"samples_per_round"`
	RoundCount               *uint32       `json:"round_count"`
	TileWidth                *uint32       `json:"tile_width"`
	TileHeight               *uint32       `json:"tile_height"`
	MaxBounces               *uint32       `json:"max_bounces"`
	UseBVH                   *bool         `json:"use_bvh"`
	BackgroundColorMagnitude *float64      `json:"background_color_magnitude"`
	BackgroundColor          *ColorRequest `json:"background_color"`
	TMin                     *float64      `json:"t_min"`
	TMax                     *float64      `json:"t_max"`
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
	getResponse := GetResponse{
		ParametersName:           parameters.ParametersName,
		ImageWidth:               parameters.ImageWidth,
		ImageHeight:              parameters.ImageHeight,
		FileType:                 string(parameters.FileType),
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
		BackgroundColor:          parameters.BackgroundColor,
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
		postRequest.BackgroundColor.Red == nil ||
		postRequest.BackgroundColor.Green == nil ||
		postRequest.BackgroundColor.Blue == nil ||
		postRequest.TMin == nil ||
		postRequest.TMax == nil {

		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// validate input
	var errorMessage = ""
	if *postRequest.ImageWidth < constants.ParametersMinimumDimension || *postRequest.ImageHeight < constants.ParametersMinimumDimension {
		errorMessage = fmt.Sprintf("image dimensions cannot be below %d in any dimension", constants.ParametersMinimumDimension)
	}
	if *postRequest.ImageWidth > constants.ParametersMaximumDimension || *postRequest.ImageHeight > constants.ParametersMaximumDimension {
		errorMessage = fmt.Sprintf("image dimensions cannot exceed %d in any dimension", constants.ParametersMinimumDimension)
	}
	if (*postRequest.ImageWidth)*(*postRequest.ImageHeight) < constants.ParametersMinimumTotalPixels {
		errorMessage = fmt.Sprintf("image size cannot be below %d pixels", constants.ParametersMinimumDimension)
	}
	if (*postRequest.ImageWidth)*(*postRequest.ImageHeight) > constants.ParametersMaximumTotalPixels {
		errorMessage = fmt.Sprintf("image size cannot exceed %d pixels", constants.ParametersMinimumDimension)
	}
	switch filetype.FileType(strings.ToUpper(*postRequest.FileType)) {
	case filetype.PNG, filetype.JPEG:
	default:
		errorMessage = "invalid file_type"
	}
	*postRequest.FileType = strings.ToUpper(*postRequest.FileType)
	if *postRequest.GammaCorrection <= 0.0 {
		errorMessage = "gamma_correction must be greater than zero"
	}
	if *postRequest.TextureGamma <= 0.0 {
		errorMessage = "texture_gamma must be greater than zero"
	}
	if *postRequest.SamplesPerRound <= 0 {
		errorMessage = "samples_per_round must be greater than zero"
	}
	if *postRequest.RoundCount <= 0 {
		errorMessage = "round_count must be greater than zero"
	}
	if *postRequest.TileWidth > *postRequest.ImageWidth || *postRequest.TileHeight > *postRequest.ImageHeight {
		errorMessage = "tile dimension must not exceed image dimension"
	}
	if *postRequest.TileWidth <= 0 || *postRequest.TileHeight <= 0 {
		errorMessage = "tile dimensions must be greater than zero"
	}
	if *postRequest.MaxBounces <= 0 {
		errorMessage = "max_bounces must be greater than zero"
	}
	if *postRequest.MaxBounces > constants.ParametersMaximumMaxBounces {
		errorMessage = fmt.Sprintf("max_bounces must not exceed %d", constants.ParametersMaximumMaxBounces)
	}
	if *postRequest.BackgroundColorMagnitude < 0.0 {
		errorMessage = "background_color_magnitude must be greater than or equal to zero"
	}
	if *postRequest.BackgroundColor.Red < 0.0 || *postRequest.BackgroundColor.Green < 0.0 || *postRequest.BackgroundColor.Blue < 0.0 {
		errorMessage = "background_color fields must be greater than or equal to zero"
	}
	if *postRequest.TMin <= 0.0 {
		errorMessage = "t_min field must be greater than zero"
	}
	if *postRequest.TMax < 0.0 {
		errorMessage = "t_max field must be greater than zero"
	}
	if *postRequest.TMax <= *postRequest.TMin {
		errorMessage = "t_max field must be greater than t_min"
	}
	if *postRequest.TMax > constants.ParametersMaximumTMax {
		errorMessage = fmt.Sprintf("t_max field must not exceed %f", constants.ParametersMaximumTMax)
	}

	// send error
	if errorMessage != "" {
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
	parameters := &parameterspersistence.Parameters{
		ParametersName:           *(postRequest.ParametersName),
		ImageWidth:               *(postRequest.ImageWidth),
		ImageHeight:              *(postRequest.ImageHeight),
		FileType:                 filetype.FileType(*(postRequest.FileType)),
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
		BackgroundColor: shading.Color{
			Red:   *(postRequest.BackgroundColor.Red),
			Green: *(postRequest.BackgroundColor.Green),
			Blue:  *(postRequest.BackgroundColor.Blue),
		},
		TMin: *(postRequest.TMin),
		TMax: *(postRequest.TMax),
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
