package primitivecontroller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/config/geometry"
	"github.com/paulwrubel/photolum/controller"
	"github.com/paulwrubel/photolum/enumeration/axis"
	"github.com/paulwrubel/photolum/enumeration/primitivetype"
	"github.com/paulwrubel/photolum/enumeration/rotationorder"
	"github.com/paulwrubel/photolum/persistence/primitivepersistence"
	"github.com/sirupsen/logrus"
)

var getEndpoint = "/primitives.GET"
var postEndpoint = "/primitives.POST"

type GetRequest struct {
	PrimitiveName *string `json:"primitive_name"`
}

type ParticipatingVolumeGetResponse struct {
	PrimitiveName             string  `json:"primitive_name"`
	PrimitiveType             string  `json:"primitive_type"`
	EncapsulatedPrimitiveName string  `json:"encapsulated_primitive_name"`
	Density                   float64 `json:"density"`
}

type SphereGetResponse struct {
	PrimitiveName      string         `json:"primitive_name"`
	PrimitiveType      string         `json:"primitive_type"`
	Center             geometry.Point `json:"center"`
	Radius             float64        `json:"radius"`
	HasInvertedNormals bool           `json:"has_inverted_normals"`
}

type CylinderGetResponse struct {
	PrimitiveName string         `json:"primitive_name"`
	PrimitiveType string         `json:"primitive_type"`
	A             geometry.Point `json:"a"`
	B             geometry.Point `json:"b"`
	Radius        float64        `json:"radius"`
}

type HollowCylinderGetResponse struct {
	PrimitiveName string         `json:"primitive_name"`
	PrimitiveType string         `json:"primitive_type"`
	A             geometry.Point `json:"a"`
	B             geometry.Point `json:"b"`
	InnerRadius   float64        `json:"inner_radius"`
	OuterRadius   float64        `json:"outer_radius"`
}

type RectangleGetResponse struct {
	PrimitiveName     string         `json:"primitive_name"`
	PrimitiveType     string         `json:"primitive_type"`
	A                 geometry.Point `json:"a"`
	B                 geometry.Point `json:"b"`
	IsCulled          bool           `json:"is_culled"`
	HasNegativeNormal bool           `json:"has_negative_normal"`
}

type TriangleGetResponse struct {
	PrimitiveName string         `json:"primitive_name"`
	PrimitiveType string         `json:"primitive_type"`
	A             geometry.Point `json:"a"`
	B             geometry.Point `json:"b"`
	C             geometry.Point `json:"c"`
}

type TriangleWithNormalsGetResponse struct {
	PrimitiveName string          `json:"primitive_name"`
	PrimitiveType string          `json:"primitive_type"`
	A             geometry.Point  `json:"a"`
	B             geometry.Point  `json:"b"`
	C             geometry.Point  `json:"c"`
	ANormal       geometry.Vector `json:"a_normal"`
	BNormal       geometry.Vector `json:"b_normal"`
	CNormal       geometry.Vector `json:"c_normal"`
}

type PlaneGetResponse struct {
	PrimitiveName string          `json:"primitive_name"`
	PrimitiveType string          `json:"primitive_type"`
	Point         geometry.Point  `json:"point"`
	Normal        geometry.Vector `json:"normal"`
	IsCulled      bool            `json:"is_culled"`
}

type PyramidGetResponse struct {
	PrimitiveName string         `json:"primitive_name"`
	PrimitiveType string         `json:"primitive_type"`
	A             geometry.Point `json:"a"`
	B             geometry.Point `json:"b"`
	Height        float64        `json:"height"`
}

type BoxGetResponse struct {
	PrimitiveName      string         `json:"primitive_name"`
	PrimitiveType      string         `json:"primitive_type"`
	A                  geometry.Point `json:"a"`
	B                  geometry.Point `json:"b"`
	HasInvertedNormals bool           `json:"has_inverted_normals"`
}

type TranslationGetResponse struct {
	PrimitiveName             string          `json:"primitive_name"`
	PrimitiveType             string          `json:"primitive_type"`
	EncapsulatedPrimitiveName string          `json:"encapsulated_primitive_name"`
	Displacement              geometry.Vector `json:"displacement"`
}

type RotationGetResponse struct {
	PrimitiveName             string  `json:"primitive_name"`
	PrimitiveType             string  `json:"primitive_type"`
	EncapsulatedPrimitiveName string  `json:"encapsulated_primitive_name"`
	Axis                      string  `json:"axis"`
	Angle                     float64 `json:"angle"`
}

type QuaternionGetResponse struct {
	PrimitiveName             string    `json:"primitive_name"`
	PrimitiveType             string    `json:"primitive_type"`
	EncapsulatedPrimitiveName string    `json:"encapsulated_primitive_name"`
	AxisAngles                []float64 `json:"axis_angles"`
	RotationOrder             string    `json:"rotation_order"`
}

type VectorRequest struct {
	X *float64 `json:"x"`
	Y *float64 `json:"y"`
	Z *float64 `json:"z"`
}

type PostRequest struct {
	PrimitiveName             *string        `json:"primitive_name"`
	PrimitiveType             *string        `json:"primitive_type"`
	EncapsulatedPrimitiveName *string        `json:"encapsulated_primitive_name"`
	A                         *VectorRequest `json:"a"`
	B                         *VectorRequest `json:"b"`
	C                         *VectorRequest `json:"c"`
	ANormal                   *VectorRequest `json:"a_normal"`
	BNormal                   *VectorRequest `json:"b_normal"`
	CNormal                   *VectorRequest `json:"c_normal"`
	Point                     *VectorRequest `json:"point"`
	Normal                    *VectorRequest `json:"normal"`
	Center                    *VectorRequest `json:"center"`
	Axis                      *string        `json:"axis"`
	Displacement              *VectorRequest `json:"displacement"`
	AxisAngles                []float64      `json:"axis_angles"`
	RotationOrder             *string        `json:"rotation_order"`
	Radius                    *float64       `json:"radius"`
	InnerRadius               *float64       `json:"inner_radius"`
	OuterRadius               *float64       `json:"outer_radius"`
	Height                    *float64       `json:"height"`
	Angle                     *float64       `json:"angle"`
	Density                   *float64       `json:"density"`
	IsCulled                  *bool          `json:"is_culled"`
	HasNegativeNormal         *bool          `json:"has_negative_normal"`
	HasInvertedNormals        *bool          `json:"has_inverted_normals"`
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
	if getRequest.PrimitiveName == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if row exists
	exists, err := primitivepersistence.DoesExist(plData, log, *getRequest.PrimitiveName)
	if err != nil {
		errorMessage := "error checking primitive existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if !exists {
		errorMessage := "primitive row does not exist"
		errorStatusCode := http.StatusNotFound

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// get from db
	primitive, err := primitivepersistence.Get(plData, log, *getRequest.PrimitiveName)
	if err != nil {
		errorMessage := "error getting primitive from database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	var getResponse interface{}
	switch primitivetype.PrimitiveType(primitive.PrimitiveType) {
	case primitivetype.ParticipatingVolume:
		getResponse = ParticipatingVolumeGetResponse{
			PrimitiveName:             primitive.PrimitiveName,
			PrimitiveType:             primitive.PrimitiveType,
			EncapsulatedPrimitiveName: *primitive.EncapsulatedPrimitiveName,
			Density:                   *primitive.Density,
		}
	case primitivetype.Sphere:
		getResponse = SphereGetResponse{
			PrimitiveName: primitive.PrimitiveName,
			PrimitiveType: primitive.PrimitiveType,
			Center: geometry.Point{
				X: primitive.Center[0],
				Y: primitive.Center[1],
				Z: primitive.Center[2],
			},
			Radius:             *primitive.Radius,
			HasInvertedNormals: *primitive.HasInvertedNormals,
		}
	case primitivetype.Cylinder:
		getResponse = CylinderGetResponse{
			PrimitiveName: primitive.PrimitiveName,
			PrimitiveType: primitive.PrimitiveType,
			A: geometry.Point{
				X: primitive.A[0],
				Y: primitive.A[1],
				Z: primitive.A[2],
			},
			B: geometry.Point{
				X: primitive.B[0],
				Y: primitive.B[1],
				Z: primitive.B[2],
			},
			Radius: *primitive.Radius,
		}
	case primitivetype.HollowCylinder:
		getResponse = HollowCylinderGetResponse{
			PrimitiveName: primitive.PrimitiveName,
			PrimitiveType: primitive.PrimitiveType,
			A: geometry.Point{
				X: primitive.A[0],
				Y: primitive.A[1],
				Z: primitive.A[2],
			},
			B: geometry.Point{
				X: primitive.B[0],
				Y: primitive.B[1],
				Z: primitive.B[2],
			},
			InnerRadius: *primitive.InnerRadius,
			OuterRadius: *primitive.OuterRadius,
		}
	case primitivetype.Rectangle:
		getResponse = RectangleGetResponse{
			PrimitiveName: primitive.PrimitiveName,
			PrimitiveType: primitive.PrimitiveType,
			A: geometry.Point{
				X: primitive.A[0],
				Y: primitive.A[1],
				Z: primitive.A[2],
			},
			B: geometry.Point{
				X: primitive.B[0],
				Y: primitive.B[1],
				Z: primitive.B[2],
			},
			IsCulled:          *primitive.IsCulled,
			HasNegativeNormal: *primitive.HasNegativeNormal,
		}
	case primitivetype.Triangle:
		if primitive.ANormal == nil || primitive.BNormal == nil || primitive.CNormal == nil {
			getResponse = TriangleGetResponse{
				PrimitiveName: primitive.PrimitiveName,
				PrimitiveType: primitive.PrimitiveType,
				A: geometry.Point{
					X: primitive.A[0],
					Y: primitive.A[1],
					Z: primitive.A[2],
				},
				B: geometry.Point{
					X: primitive.B[0],
					Y: primitive.B[1],
					Z: primitive.B[2],
				},
				C: geometry.Point{
					X: primitive.C[0],
					Y: primitive.C[1],
					Z: primitive.C[2],
				},
			}
		} else {
			getResponse = TriangleWithNormalsGetResponse{
				PrimitiveName: primitive.PrimitiveName,
				PrimitiveType: primitive.PrimitiveType,
				A: geometry.Point{
					X: primitive.A[0],
					Y: primitive.A[1],
					Z: primitive.A[2],
				},
				B: geometry.Point{
					X: primitive.B[0],
					Y: primitive.B[1],
					Z: primitive.B[2],
				},
				C: geometry.Point{
					X: primitive.C[0],
					Y: primitive.C[1],
					Z: primitive.C[2],
				},
				ANormal: geometry.Vector{
					X: primitive.ANormal[0],
					Y: primitive.ANormal[1],
					Z: primitive.ANormal[2],
				},
				BNormal: geometry.Vector{
					X: primitive.BNormal[0],
					Y: primitive.BNormal[1],
					Z: primitive.BNormal[2],
				},
				CNormal: geometry.Vector{
					X: primitive.CNormal[0],
					Y: primitive.CNormal[1],
					Z: primitive.CNormal[2],
				},
			}
		}
	case primitivetype.Plane:
		getResponse = PlaneGetResponse{
			PrimitiveName: primitive.PrimitiveName,
			PrimitiveType: primitive.PrimitiveType,
			Point: geometry.Point{
				X: primitive.Point[0],
				Y: primitive.Point[1],
				Z: primitive.Point[2],
			},
			Normal: geometry.Vector{
				X: primitive.Normal[0],
				Y: primitive.Normal[1],
				Z: primitive.Normal[2],
			},
			IsCulled: *primitive.IsCulled,
		}
	case primitivetype.Pyramid:
		getResponse = PyramidGetResponse{
			PrimitiveName: primitive.PrimitiveName,
			PrimitiveType: primitive.PrimitiveType,
			A: geometry.Point{
				X: primitive.A[0],
				Y: primitive.A[1],
				Z: primitive.A[2],
			},
			B: geometry.Point{
				X: primitive.B[0],
				Y: primitive.B[1],
				Z: primitive.B[2],
			},
			Height: *primitive.Height,
		}
	case primitivetype.Box:
		getResponse = BoxGetResponse{
			PrimitiveName: primitive.PrimitiveName,
			PrimitiveType: primitive.PrimitiveType,
			A: geometry.Point{
				X: primitive.A[0],
				Y: primitive.A[1],
				Z: primitive.A[2],
			},
			B: geometry.Point{
				X: primitive.B[0],
				Y: primitive.B[1],
				Z: primitive.B[2],
			},
			HasInvertedNormals: *primitive.HasInvertedNormals,
		}
	case primitivetype.Translation:
		getResponse = TranslationGetResponse{
			PrimitiveName:             primitive.PrimitiveName,
			PrimitiveType:             primitive.PrimitiveType,
			EncapsulatedPrimitiveName: *primitive.EncapsulatedPrimitiveName,
			Displacement: geometry.Vector{
				X: primitive.Displacement[0],
				Y: primitive.Displacement[1],
				Z: primitive.Displacement[2],
			},
		}
	case primitivetype.Rotation:
		getResponse = RotationGetResponse{
			PrimitiveName:             primitive.PrimitiveName,
			PrimitiveType:             primitive.PrimitiveType,
			EncapsulatedPrimitiveName: *primitive.EncapsulatedPrimitiveName,
			Axis:                      *primitive.Axis,
			Angle:                     *primitive.Angle,
		}
	case primitivetype.Quaternion:
		getResponse = QuaternionGetResponse{
			PrimitiveName:             primitive.PrimitiveName,
			PrimitiveType:             primitive.PrimitiveType,
			EncapsulatedPrimitiveName: *primitive.EncapsulatedPrimitiveName,
			AxisAngles:                primitive.AxisAngles,
			RotationOrder:             *primitive.RotationOrder,
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
	if postRequest.PrimitiveName == nil ||
		postRequest.PrimitiveType == nil {
		errorMessage := "missing field from request"
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// validate input
	errorMessage := ""

	// do the named encapsulated primitives exist?
	if postRequest.EncapsulatedPrimitiveName != nil {
		exists, err := primitivepersistence.DoesExist(plData, log, *postRequest.EncapsulatedPrimitiveName)
		if err != nil {
			errorMessage := "error checking primitive existence in database"
			errorStatusCode := http.StatusInternalServerError

			log.WithError(err).Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
			return
		}
		if !exists {
			errorMessage = "named encapsulated_primitive does not exist"
		}
	}

	switch primitivetype.PrimitiveType(strings.ToUpper(*postRequest.PrimitiveType)) {
	case primitivetype.ParticipatingVolume:
		if postRequest.EncapsulatedPrimitiveName == nil ||
			postRequest.Density == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		if *postRequest.Density <= 0.0 {
			errorMessage = "density must be greater than zero"
		}
	case primitivetype.Sphere:
		if postRequest.Center == nil ||
			postRequest.Radius == nil ||
			postRequest.HasInvertedNormals == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		if *postRequest.Radius <= 0.0 {
			errorMessage = "radius must be greater than zero"
		}
	case primitivetype.Cylinder:
		if postRequest.A == nil ||
			postRequest.B == nil ||
			postRequest.Radius == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		aPoint := geometry.Point{
			X: *postRequest.A.X,
			Y: *postRequest.A.Y,
			Z: *postRequest.A.Z,
		}
		bPoint := geometry.Point{
			X: *postRequest.B.X,
			Y: *postRequest.B.Y,
			Z: *postRequest.B.Z,
		}
		if *postRequest.Radius <= 0.0 {
			errorMessage = "radius must be greater than zero"
		} else if aPoint == bPoint {
			errorMessage = "a must not equal b"
		}
	case primitivetype.HollowCylinder:
		if postRequest.A == nil ||
			postRequest.B == nil ||
			postRequest.InnerRadius == nil ||
			postRequest.OuterRadius == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		aPoint := geometry.Point{
			X: *postRequest.A.X,
			Y: *postRequest.A.Y,
			Z: *postRequest.A.Z,
		}
		bPoint := geometry.Point{
			X: *postRequest.B.X,
			Y: *postRequest.B.Y,
			Z: *postRequest.B.Z,
		}
		if *postRequest.InnerRadius <= 0.0 {
			errorMessage = "inner radius must be greater than zero"
		} else if *postRequest.OuterRadius <= 0.0 {
			errorMessage = "outer radius must be greater than zero"
		} else if *postRequest.InnerRadius >= *postRequest.OuterRadius {
			errorMessage = "inner radius must not be greater than or equal to outer radius"
		} else if aPoint == bPoint {
			errorMessage = "a must not equal b"
		}
	case primitivetype.Rectangle:
		if postRequest.A == nil ||
			postRequest.B == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		equalityCheck := 0
		if *postRequest.A.X == *postRequest.B.X {
			equalityCheck++
		}
		if *postRequest.A.Y == *postRequest.B.Y {
			equalityCheck++
		}
		if *postRequest.A.Z == *postRequest.B.Z {
			equalityCheck++
		}
		switch equalityCheck {
		case 0:
			errorMessage = "points do not resolve to axis-aligned plane"
		case 1:
			// perfect! no error here
		case 2:
			errorMessage = "rectangle must not resolve to a line"
		case 3:
			errorMessage = "rectangle must not resolve to a point"
		}
	case primitivetype.Triangle:
		if postRequest.A == nil ||
			postRequest.B == nil ||
			postRequest.C == nil ||
			postRequest.IsCulled == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		aPoint := geometry.Point{
			X: *postRequest.A.X,
			Y: *postRequest.A.Y,
			Z: *postRequest.A.Z,
		}
		bPoint := geometry.Point{
			X: *postRequest.B.X,
			Y: *postRequest.B.Y,
			Z: *postRequest.B.Z,
		}
		cPoint := geometry.Point{
			X: *postRequest.C.X,
			Y: *postRequest.C.Y,
			Z: *postRequest.C.Z,
		}
		equalityCheck := 0
		if aPoint == bPoint {
			equalityCheck++
		}
		if aPoint == cPoint {
			equalityCheck++
		}
		if bPoint == cPoint {
			equalityCheck++
		}
		switch equalityCheck {
		case 0:
			// perfect! no error here
		case 1:
			errorMessage = "triangle must not resolve to a line"
		case 2:
			// how did this happen?
			errorMessage = "triangle must not break the transitive property of equality in mathematics (how did you DO that???)"
		case 3:
			errorMessage = "triangle must not resolve to a point"
		}

		normalCount := 0
		if postRequest.ANormal != nil {
			normalCount++
		}
		if postRequest.BNormal != nil {
			normalCount++
		}
		if postRequest.CNormal != nil {
			normalCount++
		}
		if normalCount != 0 && normalCount != 3 {
			errorMessage = "triangle must have either all vertex normals specified or no vertex normals specified"
		} else if normalCount == 3 {
			aNormal := geometry.Vector{
				X: *postRequest.ANormal.X,
				Y: *postRequest.ANormal.Y,
				Z: *postRequest.ANormal.Z,
			}
			bNormal := geometry.Vector{
				X: *postRequest.BNormal.X,
				Y: *postRequest.BNormal.Y,
				Z: *postRequest.BNormal.Z,
			}
			cNormal := geometry.Vector{
				X: *postRequest.CNormal.X,
				Y: *postRequest.CNormal.Y,
				Z: *postRequest.CNormal.Z,
			}
			if aNormal.Magnitude() == 0.0 {
				errorMessage = "triangle a_normal must not have zero magnitude"
			}
			if bNormal.Magnitude() == 0.0 {
				errorMessage = "triangle b_normal must not have zero magnitude"
			}
			if cNormal.Magnitude() == 0.0 {
				errorMessage = "triangle c_normal must not have zero magnitude"
			}
		}
	case primitivetype.Plane:
		if postRequest.Point == nil ||
			postRequest.Normal == nil ||
			postRequest.IsCulled == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		normal := geometry.Vector{
			X: *postRequest.Normal.X,
			Y: *postRequest.Normal.Y,
			Z: *postRequest.Normal.Z,
		}
		if normal.Magnitude() == 0.0 {
			errorMessage = "normal must not be zero vector"
		}
	case primitivetype.Pyramid:
		if postRequest.A == nil ||
			postRequest.B == nil ||
			postRequest.Height == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		if *postRequest.Height <= 0.0 {
			errorMessage = "pyramid height must be greater than zero"
		}
		if *postRequest.A.Y != *postRequest.B.Y {
			errorMessage = "pyramid must be directed upwards (base points must have same Y coordinate)"
		}
	case primitivetype.Box:
		if postRequest.A == nil ||
			postRequest.B == nil ||
			postRequest.HasInvertedNormals == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		aPoint := geometry.Point{
			X: *postRequest.A.X,
			Y: *postRequest.A.Y,
			Z: *postRequest.A.Z,
		}
		bPoint := geometry.Point{
			X: *postRequest.B.X,
			Y: *postRequest.B.Y,
			Z: *postRequest.B.Z,
		}
		c1 := geometry.MinComponents(aPoint, bPoint)
		c8 := geometry.MaxComponents(aPoint, bPoint)

		if c1.X == c8.X || c1.Y == c8.Y || c1.Z == c8.Z {
			errorMessage = "box resolves to point, line, or plane"
		}
	case primitivetype.Translation:
		if postRequest.EncapsulatedPrimitiveName == nil ||
			postRequest.Displacement == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
	case primitivetype.Rotation:
		if postRequest.EncapsulatedPrimitiveName == nil ||
			postRequest.Axis == nil ||
			postRequest.Angle == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		switch axis.Axis(strings.ToUpper(*postRequest.Axis)) {
		case axis.X:
		case axis.Y:
		case axis.Z:
		default:
			errorMessage = "invalid axis"
		}
		*postRequest.Axis = strings.ToUpper(*postRequest.Axis)
	case primitivetype.Quaternion:
		if postRequest.EncapsulatedPrimitiveName == nil ||
			postRequest.AxisAngles == nil ||
			postRequest.RotationOrder == nil {
			errorMessage := "missing field from request"
			errorStatusCode := http.StatusBadRequest

			log.Error(errorMessage)
			controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
			return
		}
		if len(postRequest.AxisAngles) != 3 {
			errorMessage = "invalid number of axis angles, should be 3"
		}
		switch rotationorder.RotationOrder(strings.ToUpper(*postRequest.RotationOrder)) {
		case rotationorder.XYX:
		case rotationorder.XYZ:
		case rotationorder.XZX:
		case rotationorder.XZY:
		case rotationorder.YXY:
		case rotationorder.YXZ:
		case rotationorder.YZX:
		case rotationorder.YZY:
		case rotationorder.ZXY:
		case rotationorder.ZXZ:
		case rotationorder.ZYX:
		case rotationorder.ZYZ:
		default:
			errorMessage = "invalid rotation_order"
		}
		*postRequest.RotationOrder = strings.ToUpper(*postRequest.RotationOrder)
	default:
		errorMessage = "invalid primitive_type"
	}

	// send error
	if errorMessage != "" {
		errorStatusCode := http.StatusBadRequest

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// check if row exists
	exists, err := primitivepersistence.DoesExist(plData, log, *postRequest.PrimitiveName)
	if err != nil {
		errorMessage := "error checking primitive existence in database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}
	if exists {
		errorMessage := "primitive row already exists"
		errorStatusCode := http.StatusConflict

		log.Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, nil)
		return
	}

	// assemble primitive
	var a []float64
	if postRequest.A == nil {
		a = nil
	} else {
		a = []float64{*postRequest.A.X, *postRequest.A.Y, *postRequest.A.Z}
	}
	var b []float64
	if postRequest.B == nil {
		b = nil
	} else {
		b = []float64{*postRequest.B.X, *postRequest.B.Y, *postRequest.B.Z}
	}
	var c []float64
	if postRequest.C == nil {
		c = nil
	} else {
		c = []float64{*postRequest.C.X, *postRequest.C.Y, *postRequest.C.Z}
	}
	var aNormal []float64
	if postRequest.ANormal == nil {
		aNormal = nil
	} else {
		aNormal = []float64{*postRequest.ANormal.X, *postRequest.ANormal.Y, *postRequest.ANormal.Z}
	}
	var bNormal []float64
	if postRequest.BNormal == nil {
		bNormal = nil
	} else {
		bNormal = []float64{*postRequest.BNormal.X, *postRequest.BNormal.Y, *postRequest.BNormal.Z}
	}
	var cNormal []float64
	if postRequest.CNormal == nil {
		cNormal = nil
	} else {
		cNormal = []float64{*postRequest.CNormal.X, *postRequest.CNormal.Y, *postRequest.CNormal.Z}
	}
	var point []float64
	if postRequest.Point == nil {
		point = nil
	} else {
		point = []float64{*postRequest.Point.X, *postRequest.Point.Y, *postRequest.Point.Z}
	}
	var normal []float64
	if postRequest.Normal == nil {
		normal = nil
	} else {
		normal = []float64{*postRequest.Normal.X, *postRequest.Normal.Y, *postRequest.Normal.Z}
	}
	var center []float64
	if postRequest.Center == nil {
		center = nil
	} else {
		center = []float64{*postRequest.Center.X, *postRequest.Center.Y, *postRequest.Center.Z}
	}
	var displacement []float64
	if postRequest.Displacement == nil {
		displacement = nil
	} else {
		displacement = []float64{*postRequest.Displacement.X, *postRequest.Displacement.Y, *postRequest.Displacement.Z}
	}
	primitive := &primitivepersistence.Primitive{
		PrimitiveName:             *postRequest.PrimitiveName,
		PrimitiveType:             strings.ToUpper(*postRequest.PrimitiveType),
		EncapsulatedPrimitiveName: postRequest.EncapsulatedPrimitiveName,
		A:                         a,
		B:                         b,
		C:                         c,
		ANormal:                   aNormal,
		BNormal:                   bNormal,
		CNormal:                   cNormal,
		Point:                     point,
		Normal:                    normal,
		Center:                    center,
		Axis:                      postRequest.Axis,
		Displacement:              displacement,
		AxisAngles:                postRequest.AxisAngles,
		RotationOrder:             postRequest.RotationOrder,
		Radius:                    postRequest.Radius,
		InnerRadius:               postRequest.InnerRadius,
		OuterRadius:               postRequest.OuterRadius,
		Height:                    postRequest.Height,
		Angle:                     postRequest.Angle,
		Density:                   postRequest.Density,
		IsCulled:                  postRequest.IsCulled,
		HasNegativeNormal:         postRequest.HasNegativeNormal,
		HasInvertedNormals:        postRequest.HasInvertedNormals,
	}

	// save to db
	err = primitivepersistence.Save(plData, log, primitive)
	if err != nil {
		errorMessage := "error saving primitive to database"
		errorStatusCode := http.StatusInternalServerError

		log.WithError(err).Error(errorMessage)
		controller.WriteErrorResponse(&response, errorStatusCode, errorMessage, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
	log.Debug("request completed")
}
