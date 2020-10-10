package rendercontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/enumeration/renderstatus"
	"github.com/paulwrubel/photolum/persistence/parameterspersistence"
	"github.com/paulwrubel/photolum/persistence/renderpersistence.go"
	"github.com/paulwrubel/photolum/persistence/scenepersistence"
	"github.com/paulwrubel/photolum/service/tracingservice"
	"github.com/sirupsen/logrus"
)

var getEndpoint = "/renders.GET"
var postEndpoint = "/renders.POST"

type GetRequest struct {
	RenderName *string `json:"render_name"`
}

type GetIncompleteResponse struct {
	RenderName             string `json:"render_name"`
	ParametersName         string `json:"parameters_name"`
	SceneName              string `json:"scene_name"`
	RenderStatus           string `json:"render_status"`
	CompletedRounds        string `json:"completed_rounds"`
	RoundProgress          string `json:"round_progress"`
	TotalProgress          string `json:"total_progress"`
	StartTime              string `json:"start_time"`
	ElapsedRuntime         string `json:"elapsed_runtime"`
	EstimatedTimeRemaining string `json:"estimated_time_remaining"`
	EstimatedEndTime       string `json:"estimated_end_time"`
}

type GetCompleteResponse struct {
	RenderName      string `json:"render_name"`
	ParametersName  string `json:"parameters_name"`
	SceneName       string `json:"scene_name"`
	RenderStatus    string `json:"render_status"`
	CompletedRounds string `json:"completed_rounds"`
	RoundProgress   string `json:"round_progress"`
	TotalProgress   string `json:"total_progress"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	TotalRuntime    string `json:"total_runtime"`
}

type PostRequest struct {
	RenderName     *string `json:"render_name"`
	ParametersName *string `json:"parameters_name"`
	SceneName      *string `json:"scene_name"`
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
	if getRequest.RenderName == nil {
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

	// get parameters from db
	parameters, err := parameterspersistence.Get(plData, log, render.ParametersName)
	if err != nil {
		errorMessage := "error getting parameters from database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	roundPercentage := 1.0 / float64(parameters.RoundCount)
	totalProgress := (float64(render.CompletedRounds) / float64(parameters.RoundCount)) + roundPercentage*(render.RoundProgress)
	var getResponse interface{}
	if renderstatus.RenderStatus(render.RenderStatus) == renderstatus.Completed {
		totalRuntime := render.EndTimestamp.Sub(render.StartTimestamp)
		getResponse = GetCompleteResponse{
			RenderName:      render.RenderName,
			ParametersName:  render.ParametersName,
			SceneName:       render.SceneName,
			RenderStatus:    render.RenderStatus,
			CompletedRounds: fmt.Sprintf("%d/%d", render.CompletedRounds, parameters.RoundCount),
			RoundProgress:   fmt.Sprintf("%.3f%%", 100*float64(render.RoundProgress)),
			TotalProgress:   fmt.Sprintf("%.3f%%", 100*totalProgress),
			StartTime:       render.StartTimestamp.Local().Format("2006-01-02 15:04:05 MST"),
			EndTime:         render.EndTimestamp.Local().Format("2006-01-02 15:04:05 MST"),
			TotalRuntime:    totalRuntime.Round(time.Second).String(),
		}
	} else {
		elapsedRuntime := time.Since(render.StartTimestamp)
		estimatedTimeRemaining := time.Duration(((1.0/totalProgress)*(1.0-totalProgress))*float64(elapsedRuntime.Nanoseconds())) * time.Nanosecond
		estimatedEndTime := render.StartTimestamp.Add(elapsedRuntime).Add(estimatedTimeRemaining)
		getResponse = GetIncompleteResponse{
			RenderName:             render.RenderName,
			ParametersName:         render.ParametersName,
			SceneName:              render.SceneName,
			RenderStatus:           render.RenderStatus,
			CompletedRounds:        fmt.Sprintf("%d/%d", render.CompletedRounds, parameters.RoundCount),
			RoundProgress:          fmt.Sprintf("%.3f%%", 100*float64(render.RoundProgress)),
			TotalProgress:          fmt.Sprintf("%.3f%%", 100*totalProgress),
			StartTime:              render.StartTimestamp.Local().Format("2006-01-02 15:04:05 MST"),
			ElapsedRuntime:         elapsedRuntime.String(),
			EstimatedTimeRemaining: estimatedTimeRemaining.String(),
			EstimatedEndTime:       estimatedEndTime.Local().Format("2006-01-02 15:04:05 MST"),
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
	if postRequest.RenderName == nil ||
		postRequest.ParametersName == nil ||
		postRequest.SceneName == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	errorMessage := ""
	// check if parameters exists
	exists, err := parameterspersistence.DoesExist(plData, log, *postRequest.ParametersName)
	if err != nil {
		errorMessage := "error checking parameters existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if !exists {
		errorMessage = "named parameters does not exist"
	}
	// check if scene exists
	exists, err = scenepersistence.DoesExist(plData, log, *postRequest.SceneName)
	if err != nil {
		errorMessage := "error checking scene existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if !exists {
		errorMessage = "named scene does not exist"
	}

	// send error
	if errorMessage != "" {
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if render row exists
	exists, err = renderpersistence.DoesExist(plData, log, *postRequest.RenderName)
	if err != nil {
		errorMessage := "error checking render existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if exists {
		errorMessage := "render row already exists"
		errorStatusCode := http.StatusConflict

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// assemble render
	render := &renderpersistence.Render{
		RenderName:      *postRequest.RenderName,
		ParametersName:  *postRequest.ParametersName,
		SceneName:       *postRequest.SceneName,
		RenderStatus:    string(renderstatus.Created),
		CompletedRounds: 0,
		RoundProgress:   0.0,
		StartTimestamp:  time.Now(),
	}

	// save render to db
	err = renderpersistence.Save(plData, log, render)
	if err != nil {
		errorMessage := "error saving render to database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	// set render to pending
	renderpersistence.UpdateRenderStatus(plData, log, render.RenderName, renderstatus.Pending)
	if err != nil {
		errorMessage := "error updating render status to PENDING"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	// begin assembly and rendering
	err = tracingservice.StartRender(plData, baseLog, render.RenderName)
	if err != nil {
		errorMessage := "error starting render"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
	log.Debug("request completed")
}
