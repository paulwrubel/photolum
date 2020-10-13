package materialcontroller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/enumeration/materialtype"
	"github.com/paulwrubel/photolum/persistence/materialpersistence"
	"github.com/paulwrubel/photolum/persistence/texturepersistence"
	"github.com/sirupsen/logrus"
)

var getEndpoint = "/materials.GET"
var postEndpoint = "/materials.POST"

type GetRequest struct {
	MaterialName *string `json:"material_name"`
}

type LambertianGetResponse struct {
	MaterialName           string `json:"material_name"`
	MaterialType           string `json:"material_type"`
	ReflectanceTextureName string `json:"reflectance_texture_name,omitempty"`
	EmittanceTextureName   string `json:"emittance_texture_name,omitempty"`
}

type MetalGetResponse struct {
	MaterialName           string  `json:"material_name"`
	MaterialType           string  `json:"material_type"`
	ReflectanceTextureName string  `json:"reflectance_texture_name,omitempty"`
	EmittanceTextureName   string  `json:"emittance_texture_name,omitempty"`
	Fuzziness              float64 `json:"fuzziness"`
}

type DielectricGetResponse struct {
	MaterialName           string  `json:"material_name"`
	MaterialType           string  `json:"material_type"`
	ReflectanceTextureName string  `json:"reflectance_texture_name,omitempty"`
	EmittanceTextureName   string  `json:"emittance_texture_name,omitempty"`
	RefractiveIndex        float64 `json:"refractive_index"`
}

type IsotropicGetResponse struct {
	MaterialName           string `json:"material_name"`
	MaterialType           string `json:"material_type"`
	ReflectanceTextureName string `json:"reflectance_texture_name,omitempty"`
	EmittanceTextureName   string `json:"emittance_texture_name,omitempty"`
}

type PostRequest struct {
	MaterialName           *string  `json:"material_name"`
	MaterialType           *string  `json:"material_type"`
	ReflectanceTextureName *string  `json:"reflectance_texture_name"`
	EmittanceTextureName   *string  `json:"emittance_texture_name"`
	Fuzziness              *float64 `json:"fuzziness"`
	RefractiveIndex        *float64 `json:"refractive_index"`
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
	if getRequest.MaterialName == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if row exists
	exists, err := materialpersistence.DoesExist(plData, log, *getRequest.MaterialName)
	if err != nil {
		errorMessage := "error checking material existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if !exists {
		errorMessage := "material row does not exist"
		errorStatusCode := http.StatusNotFound

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// get from db
	material, err := materialpersistence.Get(plData, log, *getRequest.MaterialName)
	if err != nil {
		errorMessage := "error getting material from database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	var reflectanceTextureName, emittanceTextureName string
	if material.ReflectanceTextureName != nil {
		reflectanceTextureName = *material.ReflectanceTextureName
	}
	if material.EmittanceTextureName != nil {
		emittanceTextureName = *material.EmittanceTextureName
	}
	var getResponse interface{}
	switch materialtype.MaterialType(material.MaterialType) {
	case materialtype.Lambertian:
		getResponse = LambertianGetResponse{
			MaterialName:           material.MaterialName,
			MaterialType:           material.MaterialType,
			ReflectanceTextureName: reflectanceTextureName,
			EmittanceTextureName:   emittanceTextureName,
		}
	case materialtype.Metal:
		getResponse = MetalGetResponse{
			MaterialName:           material.MaterialName,
			MaterialType:           material.MaterialType,
			ReflectanceTextureName: reflectanceTextureName,
			EmittanceTextureName:   emittanceTextureName,
			Fuzziness:              *material.Fuzziness,
		}
	case materialtype.Dielectric:
		getResponse = DielectricGetResponse{
			MaterialName:           material.MaterialName,
			MaterialType:           material.MaterialType,
			ReflectanceTextureName: reflectanceTextureName,
			EmittanceTextureName:   emittanceTextureName,
			RefractiveIndex:        *material.RefractiveIndex,
		}
	case materialtype.Isotropic:
		getResponse = IsotropicGetResponse{
			MaterialName:           material.MaterialName,
			MaterialType:           material.MaterialType,
			ReflectanceTextureName: reflectanceTextureName,
			EmittanceTextureName:   emittanceTextureName,
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
	if postRequest.MaterialName == nil ||
		postRequest.MaterialType == nil ||
		(postRequest.ReflectanceTextureName == nil && postRequest.EmittanceTextureName == nil) {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// validate input
	errorMessage := ""

	// do the named textures exist?
	if postRequest.ReflectanceTextureName != nil {
		exists, err := texturepersistence.DoesExist(plData, log, *postRequest.ReflectanceTextureName)
		if err != nil {
			errorMessage := "error checking texture existence in database"
			errorStatusCode := http.StatusInternalServerError

			log.WithError(err).Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
			return
		}
		if !exists {
			errorMessage = "named reflectance_texture does not exist"
		}
	}
	if postRequest.EmittanceTextureName != nil {
		exists, err := texturepersistence.DoesExist(plData, log, *postRequest.EmittanceTextureName)
		if err != nil {
			errorMessage := "error checking texture existence in database"
			errorStatusCode := http.StatusInternalServerError

			log.WithError(err).Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
			return
		}
		if !exists {
			errorMessage = "named emittance_texture does not exist"
		}
	}
	switch materialtype.MaterialType(strings.ToUpper(*postRequest.MaterialType)) {
	case materialtype.Lambertian:
		// no unique validation necessary
	case materialtype.Metal:
		if postRequest.Fuzziness == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		if *postRequest.Fuzziness < 0.0 {
			errorMessage = "fuzziness must be greater than or equal to zero"
		} else if *postRequest.Fuzziness >= 1.0 {
			errorMessage = "fuzziness must be less than 1.0"
		}
	case materialtype.Dielectric:
		if postRequest.RefractiveIndex == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		if *postRequest.RefractiveIndex <= 1.0 {
			errorMessage = "refractive_index must be greater than 1.0"
		}
	case materialtype.Isotropic:
		// no unique validation necessary
	default:
		errorMessage = "invalid material_type"
	}

	// send error
	if errorMessage != "" {
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if row exists
	exists, err := materialpersistence.DoesExist(plData, log, *postRequest.MaterialName)
	if err != nil {
		errorMessage := "error checking material existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if exists {
		errorMessage := "material row already exists"
		errorStatusCode := http.StatusConflict

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// assemble material
	material := &materialpersistence.Material{
		MaterialName:           *postRequest.MaterialName,
		MaterialType:           strings.ToUpper(*postRequest.MaterialType),
		ReflectanceTextureName: postRequest.ReflectanceTextureName,
		EmittanceTextureName:   postRequest.EmittanceTextureName,
		Fuzziness:              postRequest.Fuzziness,
		RefractiveIndex:        postRequest.RefractiveIndex,
	}

	// save to db
	err = materialpersistence.Save(plData, log, material)
	if err != nil {
		errorMessage := "error saving material to database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
	log.Debug("request completed")
}
