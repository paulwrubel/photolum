package tracingservice

import (
	"fmt"
	"reflect"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/config/geometry"
	"github.com/paulwrubel/photolum/config/geometry/primitive"
	"github.com/paulwrubel/photolum/config/geometry/primitive/box"
	"github.com/paulwrubel/photolum/config/geometry/primitive/bvh"
	"github.com/paulwrubel/photolum/config/geometry/primitive/cylinder"
	"github.com/paulwrubel/photolum/config/geometry/primitive/hollowcylinder"
	"github.com/paulwrubel/photolum/config/geometry/primitive/participatingvolume"
	"github.com/paulwrubel/photolum/config/geometry/primitive/plane"
	"github.com/paulwrubel/photolum/config/geometry/primitive/primitivelist"
	"github.com/paulwrubel/photolum/config/geometry/primitive/pyramid"
	"github.com/paulwrubel/photolum/config/geometry/primitive/rectangle"
	"github.com/paulwrubel/photolum/config/geometry/primitive/sphere"
	"github.com/paulwrubel/photolum/config/geometry/primitive/transform/rotate"
	"github.com/paulwrubel/photolum/config/geometry/primitive/transform/translate"
	"github.com/paulwrubel/photolum/config/geometry/primitive/triangle"
	"github.com/paulwrubel/photolum/config/shading"
	"github.com/paulwrubel/photolum/config/shading/material"
	"github.com/paulwrubel/photolum/config/shading/texture"
	"github.com/paulwrubel/photolum/encoding"
	"github.com/paulwrubel/photolum/enumeration/axis"
	"github.com/paulwrubel/photolum/enumeration/filetype"
	"github.com/paulwrubel/photolum/enumeration/materialtype"
	"github.com/paulwrubel/photolum/enumeration/primitivetype"
	"github.com/paulwrubel/photolum/enumeration/renderstatus"
	"github.com/paulwrubel/photolum/enumeration/texturetype"
	"github.com/paulwrubel/photolum/persistence/camerapersistence"
	"github.com/paulwrubel/photolum/persistence/materialpersistence"
	"github.com/paulwrubel/photolum/persistence/parameterspersistence"
	"github.com/paulwrubel/photolum/persistence/primitivepersistence"
	"github.com/paulwrubel/photolum/persistence/renderpersistence.go"
	"github.com/paulwrubel/photolum/persistence/scenepersistence"
	"github.com/paulwrubel/photolum/persistence/sceneprimitivematerialpersistence"
	"github.com/paulwrubel/photolum/persistence/texturepersistence"
	"github.com/paulwrubel/photolum/tracing"
	"github.com/sirupsen/logrus"
)

func StartRender(plData *config.PhotolumData, baseLog *logrus.Logger, renderName string) error {
	log := baseLog.WithFields(logrus.Fields{
		"render_name": renderName,
	})
	log.Debug("starting render")

	err := renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Starting)
	if err != nil {
		log.WithError(err).Error("error setting render to starting")
		renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
		return err
	}

	parameters, err := loadConfigurations(plData, log, renderName)
	if err != nil {
		log.WithError(err).Error("error loading Parameters")
		renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
		return err
	}

	//spew.Dump(parameters)

	encodingChan := make(chan *config.TracingPayload)
	// start encoding worker
	go encoding.RunWorker(plData, log, renderName, encodingChan)
	// start tracing worker
	go tracing.RunWorker(plData, log, parameters, renderName, encodingChan)

	err = renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Running)
	if err != nil {
		log.WithError(err).Error("error setting render to running")
		renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
		return err
	}

	return nil
}

func loadConfigurations(plData *config.PhotolumData, log *logrus.Entry, renderName string) (*config.Parameters, error) {
	log.Debug("loading configurations")
	// get render from db
	renderDB, err := renderpersistence.Get(plData, log, renderName)
	if err != nil {
		renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
		return nil, fmt.Errorf("error getting render from db: %s", err.Error())
	}
	// get parameters from db
	parametersDB, err := parameterspersistence.Get(plData, log, renderDB.ParametersName)
	if err != nil {
		renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
		return nil, fmt.Errorf("error getting parameters from db: %s", err.Error())
	}
	// create Parameters struct
	parameters := decodeParameters(parametersDB)

	// get scene from db
	sceneDB, err := scenepersistence.Get(plData, log, renderDB.SceneName)
	if err != nil {
		renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
		return nil, fmt.Errorf("error getting scene from db: %s", err.Error())
	}
	// create and attach scene
	parameters.Scene = &config.Scene{}

	// get camera from db
	cameraDB, err := camerapersistence.Get(plData, log, sceneDB.CameraName)
	if err != nil {
		renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
		return nil, fmt.Errorf("error getting camera from db: %s", err.Error())
	}
	// attach camera to scene
	parameters.Scene.Camera = decodeCamera(cameraDB, parameters)

	// get sceneprimitivematerials from db
	spmListDB, err := sceneprimitivematerialpersistence.GetAllInScene(plData, log, sceneDB.SceneName)
	if err != nil {
		renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
		return nil, fmt.Errorf("error getting sceneprimitivematerials from db: %s", err.Error())
	}

	// start setup attachment process

	// most geometry primitives are "bounded" meaning an AABB (Axis-Aligned Bounding Box) can be placed around them.
	boundedSceneObjects := &primitivelist.PrimitiveList{}
	// some geometry primitives, however, are infinite in nature, which mean they cannot be bounded.
	// a distinction must be made between these to prevent assembling a BVH or other acceleration structure
	// without a bounding box around certain primitives
	unboundedSceneObjects := &primitivelist.PrimitiveList{}
	for _, spm := range spmListDB {
		// get primitive from DB
		primitiveDB, err := primitivepersistence.Get(plData, log, spm.PrimitiveName)
		if err != nil {
			renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
			return nil, fmt.Errorf("error getting primitive from db: %s", err.Error())
		}
		// decode primitive
		selectedPrimitive, err := decodePrimitive(plData, log, primitiveDB)
		if err != nil {
			renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
			return nil, fmt.Errorf("error decoding primitive: %s", err.Error())
		}

		// get material from DB
		materialDB, err := materialpersistence.Get(plData, log, spm.MaterialName)
		if err != nil {
			renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
			return nil, fmt.Errorf("error getting material from db: %s", err.Error())
		}
		// decode material
		selectedMaterial, err := decodeMaterial(plData, log, materialDB)
		if err != nil {
			renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
			return nil, fmt.Errorf("error decoding material: %s", err.Error())
		}

		// this is a check to ensure that materials that have a transmission component (i.e. Dielectrics, isotropics)
		// are not attached to "open" geometry, such as single-sided triangles and rectangles, so the
		// transmission commponent can be reversed
		// this is an arbitrary restriction that is likely to be removed in the future with the user choosing to self-restrict
		// themselves in a similar manner
		if reflect.TypeOf(selectedMaterial) == reflect.TypeOf(&material.Dielectric{}) && !selectedPrimitive.IsClosed() {
			return nil, fmt.Errorf("cannot attach refractive materials (%s) to non-closed geometry (%s)",
				spm.MaterialName, spm.PrimitiveName)
		}
		if reflect.TypeOf(selectedMaterial) == reflect.TypeOf(&material.Isotropic{}) && !selectedPrimitive.IsClosed() {
			return nil, fmt.Errorf("cannot attach volumetric materials (%s) to non-closed geometry (%s)",
				spm.MaterialName, spm.PrimitiveName)
		}

		// additionally, isotropics specifically must be attached to participating volumes
		if reflect.TypeOf(selectedMaterial) == reflect.TypeOf(&material.Isotropic{}) &&
			reflect.TypeOf(selectedPrimitive) != reflect.TypeOf(&participatingvolume.ParticipatingVolume{}) {
			return nil, fmt.Errorf("cannot attach isotropic materials (%s) to primitive not of type participating_volume (%s)",
				spm.MaterialName, spm.PrimitiveName)
		}
		// ...and vice versa as well
		if reflect.TypeOf(selectedPrimitive) == reflect.TypeOf(&participatingvolume.ParticipatingVolume{}) &&
			reflect.TypeOf(selectedMaterial) != reflect.TypeOf(&material.Isotropic{}) {
			return nil, fmt.Errorf("cannot attach to participating volume (%s) a material not of type isotropic (%s)",
				spm.PrimitiveName, spm.MaterialName)
		}

		selectedPrimitive.SetMaterial(selectedMaterial)
		// added to the cooresponding list based on type
		if selectedPrimitive.IsInfinite() {
			unboundedSceneObjects.List = append(unboundedSceneObjects.List, selectedPrimitive)
		} else {
			boundedSceneObjects.List = append(boundedSceneObjects.List, selectedPrimitive)
		}
	}

	// if we are using a BVH ...
	if parameters.UseBVH {
		// ... construct it from the bounded objects ..
		sceneBVH, err := bvh.New(boundedSceneObjects)
		if err != nil {
			return nil, err
		}
		// ... and set it as the root node if no infinite geometry exists
		if len(unboundedSceneObjects.List) == 0 {
			parameters.Scene.Objects = sceneBVH
		} else {
			// but if some infinite geometry exists in the scene, we then
			// establish a new root node as a list of the BVH and the infinite geometry
			rootNode := &primitivelist.PrimitiveList{
				List: append(unboundedSceneObjects.List, sceneBVH),
			}
			parameters.Scene.Objects = rootNode
		}
	} else {
		// if we are not using a BVH, combine the lists into a core list and set it as the root node
		parameters.Scene.Objects = &primitivelist.PrimitiveList{
			List: append(boundedSceneObjects.List, unboundedSceneObjects.List...),
		}
	}

	return parameters, nil
}

func decodeParameters(parametersDB *parameterspersistence.Parameters) *config.Parameters {
	parameters := &config.Parameters{
		ImageWidth:               int(parametersDB.ImageWidth),
		ImageHeight:              int(parametersDB.ImageHeight),
		FileType:                 filetype.FileType(parametersDB.FileType),
		GammaCorrection:          parametersDB.GammaCorrection,
		UseScalingTruncation:     parametersDB.UseScalingTruncation,
		SamplesPerRound:          int(parametersDB.SamplesPerRound),
		RoundCount:               int(parametersDB.RoundCount),
		TileWidth:                int(parametersDB.TileWidth),
		TileHeight:               int(parametersDB.TileHeight),
		MaxBounces:               int(parametersDB.MaxBounces),
		UseBVH:                   parametersDB.UseBVH,
		BackgroundColorMagnitude: parametersDB.BackgroundColorMagnitude,
		BackgroundColor: shading.Color{
			Red:   parametersDB.BackgroundColor[0],
			Green: parametersDB.BackgroundColor[1],
			Blue:  parametersDB.BackgroundColor[2],
		},
		TMin: parametersDB.TMin,
		TMax: parametersDB.TMax,
	}
	parameters.BackgroundColor = parameters.BackgroundColor.MultScalar(parameters.BackgroundColorMagnitude)
	return parameters
}

func decodeCamera(cameraDB *camerapersistence.Camera, parameters *config.Parameters) *config.Camera {
	camera := &config.Camera{
		EyeLocation: geometry.Point{
			X: cameraDB.EyeLocation[0],
			Y: cameraDB.EyeLocation[1],
			Z: cameraDB.EyeLocation[2],
		},
		TargetLocation: geometry.Point{
			X: cameraDB.TargetLocation[0],
			Y: cameraDB.TargetLocation[1],
			Z: cameraDB.TargetLocation[2],
		},
		UpVector: geometry.Vector{
			X: cameraDB.UpVector[0],
			Y: cameraDB.UpVector[1],
			Z: cameraDB.UpVector[2],
		},
		VerticalFOV:   cameraDB.VerticalFOV,
		Aperture:      cameraDB.Aperture,
		FocusDistance: cameraDB.FocusDistance,
	}
	camera.Setup(parameters)
	return camera
}

func decodePrimitive(plData *config.PhotolumData, log *logrus.Entry, primitiveDB *primitivepersistence.Primitive) (primitive.Primitive, error) {
	switch primitivetype.PrimitiveType(primitiveDB.PrimitiveType) {
	case primitivetype.ParticipatingVolume:
		corePrimitiveDB, err := primitivepersistence.Get(plData, log, *primitiveDB.EncapsulatedPrimitiveName)
		if err != nil {
			return nil, err
		}
		corePrimitive, err := decodePrimitive(plData, log, corePrimitiveDB)
		if err != nil {
			return nil, err
		}
		newPV, err := (&participatingvolume.ParticipatingVolume{
			Density:   *primitiveDB.Density,
			Primitive: corePrimitive,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newPV, nil
	case primitivetype.Box:
		newBox, err := (&box.Box{
			A: geometry.Point{
				X: primitiveDB.A[0],
				Y: primitiveDB.A[1],
				Z: primitiveDB.A[2],
			},
			B: geometry.Point{
				X: primitiveDB.B[0],
				Y: primitiveDB.B[1],
				Z: primitiveDB.B[2],
			},
			HasInvertedNormals: *primitiveDB.HasInvertedNormals,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newBox, nil
	case primitivetype.Cylinder:
		newCylinder, err := (&cylinder.Cylinder{
			A: geometry.Point{
				X: primitiveDB.A[0],
				Y: primitiveDB.A[1],
				Z: primitiveDB.A[2],
			},
			B: geometry.Point{
				X: primitiveDB.B[0],
				Y: primitiveDB.B[1],
				Z: primitiveDB.B[2],
			},
			Radius: *primitiveDB.Radius,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newCylinder, nil
	case primitivetype.HollowCylinder:
		newHollowCylinder, err := (&hollowcylinder.HollowCylinder{
			A: geometry.Point{
				X: primitiveDB.A[0],
				Y: primitiveDB.A[1],
				Z: primitiveDB.A[2],
			},
			B: geometry.Point{
				X: primitiveDB.B[0],
				Y: primitiveDB.B[1],
				Z: primitiveDB.B[2],
			},
			InnerRadius: *primitiveDB.InnerRadius,
			OuterRadius: *primitiveDB.OuterRadius,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newHollowCylinder, nil
	// case "InfiniteCylinder":
	// case "UncappedCylinder":
	// case "Disk":
	// case "HollowDisk":
	case primitivetype.Plane:
		newPlane, err := (&plane.Plane{
			Point: geometry.Point{
				X: primitiveDB.Point[0],
				Y: primitiveDB.Point[1],
				Z: primitiveDB.Point[2],
			},
			Normal: geometry.Vector{
				X: primitiveDB.Normal[0],
				Y: primitiveDB.Normal[1],
				Z: primitiveDB.Normal[2],
			},
			IsCulled: *primitiveDB.IsCulled,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newPlane, nil
	case primitivetype.Pyramid:
		newPlane, err := (&pyramid.Pyramid{
			A: geometry.Point{
				X: primitiveDB.A[0],
				Y: primitiveDB.A[1],
				Z: primitiveDB.A[2],
			},
			B: geometry.Point{
				X: primitiveDB.B[0],
				Y: primitiveDB.B[1],
				Z: primitiveDB.B[2],
			},
			Height: *primitiveDB.Height,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newPlane, nil
	case primitivetype.Rectangle:
		newRectangle, err := (&rectangle.Rectangle{
			A: geometry.Point{
				X: primitiveDB.A[0],
				Y: primitiveDB.A[1],
				Z: primitiveDB.A[2],
			},
			B: geometry.Point{
				X: primitiveDB.B[0],
				Y: primitiveDB.B[1],
				Z: primitiveDB.B[2],
			},
			IsCulled:          *primitiveDB.IsCulled,
			HasNegativeNormal: *primitiveDB.HasNegativeNormal,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newRectangle, nil
	case primitivetype.Sphere:
		newSphere, err := (&sphere.Sphere{
			Center: geometry.Point{
				X: primitiveDB.Center[0],
				Y: primitiveDB.Center[1],
				Z: primitiveDB.Center[2],
			},
			Radius:             *primitiveDB.Radius,
			HasInvertedNormals: *primitiveDB.HasInvertedNormals,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newSphere, nil
	case primitivetype.Triangle:
		var aNormal, bNormal, cNormal geometry.Vector
		if primitiveDB.ANormal != nil {
			aNormal = geometry.Vector{
				X: primitiveDB.ANormal[0],
				Y: primitiveDB.ANormal[1],
				Z: primitiveDB.ANormal[2],
			}
		}
		if primitiveDB.BNormal != nil {
			bNormal = geometry.Vector{
				X: primitiveDB.BNormal[0],
				Y: primitiveDB.BNormal[1],
				Z: primitiveDB.BNormal[2],
			}
		}
		if primitiveDB.CNormal != nil {
			cNormal = geometry.Vector{
				X: primitiveDB.CNormal[0],
				Y: primitiveDB.CNormal[1],
				Z: primitiveDB.CNormal[2],
			}
		}
		newTriangle, err := (&triangle.Triangle{
			A: geometry.Point{
				X: primitiveDB.A[0],
				Y: primitiveDB.A[1],
				Z: primitiveDB.A[2],
			},
			B: geometry.Point{
				X: primitiveDB.B[0],
				Y: primitiveDB.B[1],
				Z: primitiveDB.B[2],
			},
			C: geometry.Point{
				X: primitiveDB.C[0],
				Y: primitiveDB.C[1],
				Z: primitiveDB.C[2],
			},
			ANormal:  aNormal,
			BNormal:  bNormal,
			CNormal:  cNormal,
			IsCulled: *primitiveDB.IsCulled,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newTriangle, nil
	case primitivetype.Translation:
		corePrimitiveDB, err := primitivepersistence.Get(plData, log, *primitiveDB.EncapsulatedPrimitiveName)
		if err != nil {
			return nil, err
		}
		corePrimitive, err := decodePrimitive(plData, log, corePrimitiveDB)
		if err != nil {
			return nil, err
		}
		newTranslation, err := (&translate.Translation{
			Displacement: geometry.Vector{
				X: primitiveDB.Displacement[0],
				Y: primitiveDB.Displacement[1],
				Z: primitiveDB.Displacement[2],
			},
			Primitive: corePrimitive,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newTranslation, nil
	case primitivetype.Rotation:
		corePrimitiveDB, err := primitivepersistence.Get(plData, log, *primitiveDB.EncapsulatedPrimitiveName)
		if err != nil {
			return nil, err
		}
		corePrimitive, err := decodePrimitive(plData, log, corePrimitiveDB)
		if err != nil {
			return nil, err
		}
		var newRotation primitive.Primitive
		switch axis.Axis(*primitiveDB.Axis) {
		case axis.X:
			newRotation, err = (&rotate.RotationX{
				AngleDegrees: *primitiveDB.Angle,
				Primitive:    corePrimitive,
			}).Setup()
		case axis.Y:
			newRotation, err = (&rotate.RotationY{
				AngleDegrees: *primitiveDB.Angle,
				Primitive:    corePrimitive,
			}).Setup()
		case axis.Z:
			newRotation, err = (&rotate.RotationZ{
				AngleDegrees: *primitiveDB.Angle,
				Primitive:    corePrimitive,
			}).Setup()
		}
		if err != nil {
			return nil, err
		}
		return newRotation, nil
	case primitivetype.Quaternion:
		corePrimitiveDB, err := primitivepersistence.Get(plData, log, *primitiveDB.EncapsulatedPrimitiveName)
		if err != nil {
			return nil, err
		}
		corePrimitive, err := decodePrimitive(plData, log, corePrimitiveDB)
		if err != nil {
			return nil, err
		}
		newQuaternion, err := (&rotate.Quaternion{
			AxisAngles: primitiveDB.AxisAngles,
			Order:      *primitiveDB.RotationOrder,
			Primitive:  corePrimitive,
		}).Setup()
		if err != nil {
			return nil, err
		}
		return newQuaternion, nil
	default:
		return nil, fmt.Errorf("invalid primitive type")
	}
}

func decodeMaterial(plData *config.PhotolumData, log *logrus.Entry, materialDB *materialpersistence.Material) (material.Material, error) {
	switch materialtype.MaterialType(materialDB.MaterialType) {
	case materialtype.Lambertian:
		newMaterial := &material.Lambertian{
			ReflectanceTexture: &texture.Color{Color: shading.ColorBlack},
			EmittanceTexture:   &texture.Color{Color: shading.ColorBlack},
		}
		if materialDB.ReflectanceTextureName != nil {
			textureDB, err := texturepersistence.Get(plData, log, *materialDB.ReflectanceTextureName)
			if err != nil {
				return nil, err
			}
			newMaterial.ReflectanceTexture, err = decodeTexture(plData, log, textureDB)
			if err != nil {
				return nil, err
			}
		}
		if materialDB.EmittanceTextureName != nil {
			textureDB, err := texturepersistence.Get(plData, log, *materialDB.EmittanceTextureName)
			if err != nil {
				return nil, err
			}
			newMaterial.EmittanceTexture, err = decodeTexture(plData, log, textureDB)
			if err != nil {
				return nil, err
			}
		}
		return newMaterial, nil
	case materialtype.Metal:
		newMaterial := &material.Metal{
			ReflectanceTexture: &texture.Color{Color: shading.ColorBlack},
			EmittanceTexture:   &texture.Color{Color: shading.ColorBlack},
			Fuzziness:          *materialDB.Fuzziness,
		}
		if materialDB.ReflectanceTextureName != nil {
			textureDB, err := texturepersistence.Get(plData, log, *materialDB.ReflectanceTextureName)
			if err != nil {
				return nil, err
			}
			newMaterial.ReflectanceTexture, err = decodeTexture(plData, log, textureDB)
			if err != nil {
				return nil, err
			}
		}
		if materialDB.EmittanceTextureName != nil {
			textureDB, err := texturepersistence.Get(plData, log, *materialDB.EmittanceTextureName)
			if err != nil {
				return nil, err
			}
			newMaterial.EmittanceTexture, err = decodeTexture(plData, log, textureDB)
			if err != nil {
				return nil, err
			}
		}
		return newMaterial, nil
	case materialtype.Dielectric:
		newMaterial := &material.Dielectric{
			ReflectanceTexture: &texture.Color{Color: shading.ColorBlack},
			EmittanceTexture:   &texture.Color{Color: shading.ColorBlack},
			RefractiveIndex:    *materialDB.RefractiveIndex,
		}
		if materialDB.ReflectanceTextureName != nil {
			textureDB, err := texturepersistence.Get(plData, log, *materialDB.ReflectanceTextureName)
			if err != nil {
				return nil, err
			}
			newMaterial.ReflectanceTexture, err = decodeTexture(plData, log, textureDB)
			if err != nil {
				return nil, err
			}
		}
		if materialDB.EmittanceTextureName != nil {
			textureDB, err := texturepersistence.Get(plData, log, *materialDB.EmittanceTextureName)
			if err != nil {
				return nil, err
			}
			newMaterial.EmittanceTexture, err = decodeTexture(plData, log, textureDB)
			if err != nil {
				return nil, err
			}
		}
		return newMaterial, nil
	case materialtype.Isotropic:
		newMaterial := &material.Isotropic{
			ReflectanceTexture: &texture.Color{Color: shading.ColorBlack},
			EmittanceTexture:   &texture.Color{Color: shading.ColorBlack},
		}
		if materialDB.ReflectanceTextureName != nil {
			textureDB, err := texturepersistence.Get(plData, log, *materialDB.ReflectanceTextureName)
			if err != nil {
				return nil, err
			}
			newMaterial.ReflectanceTexture, err = decodeTexture(plData, log, textureDB)
			if err != nil {
				return nil, err
			}
		}
		if materialDB.EmittanceTextureName != nil {
			textureDB, err := texturepersistence.Get(plData, log, *materialDB.EmittanceTextureName)
			if err != nil {
				return nil, err
			}
			newMaterial.EmittanceTexture, err = decodeTexture(plData, log, textureDB)
			if err != nil {
				return nil, err
			}
		}
		return newMaterial, nil
	default:
		return nil, fmt.Errorf("invalid material type")
	}
}

func decodeTexture(plData *config.PhotolumData, log *logrus.Entry, textureDB *texturepersistence.Texture) (texture.Texture, error) {
	switch texturetype.TextureType(textureDB.TextureType) {
	case texturetype.Color:
		newTexture := &texture.Color{
			Color: shading.Color{
				Red:   textureDB.Color[0],
				Green: textureDB.Color[1],
				Blue:  textureDB.Color[2],
			},
		}
		return newTexture, nil
	case texturetype.Image:
		newTexture := &texture.Image{
			ImageData: textureDB.ImageData,
			Gamma:     *textureDB.Gamma,
			Magnitude: *textureDB.Magnitude,
		}
		err := newTexture.Load()
		if err != nil {
			return nil, err
		}
		return newTexture, nil
	default:
		return nil, fmt.Errorf("invalid texture type")
	}
}
