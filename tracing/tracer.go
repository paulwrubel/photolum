package tracing

import (
	"fmt"
	_image "image"
	"image/color"
	"math"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/config/renderstatus"
	"github.com/paulwrubel/photolum/persistence/image"
	"github.com/paulwrubel/photolum/persistence/scene"
)

func StartRender(plData *config.PhotolumData, sceneID string) {
	go TraceImage(plData, sceneID)
	scene.UpdateRenderStatus(plData, sceneID, renderstatus.Running)
}

func StopRender(plData *config.PhotolumData, sceneID string) {
	scene.UpdateRenderStatus(plData, sceneID, renderstatus.Stopping)
}

func TraceImage(plData *config.PhotolumData, sceneID string) {
	currentScene, err := scene.Retrieve(plData, sceneID)
	if err != nil {
		fmt.Printf("Error in tracing.go: %s\n", err.Error())
		scene.UpdateRenderStatus(plData, sceneID, renderstatus.Error)
		return
	}
	newImage := _image.NewRGBA64(_image.Rect(0, 0, currentScene.ImageWidth, currentScene.ImageHeight))
	for y := 0; y < currentScene.ImageHeight; y++ {
		for x := 0; x < currentScene.ImageWidth; x++ {
			col := color.RGBA64{
				R: uint16(0.0 * float64(math.MaxUint16)),
				G: uint16((float64(x) / float64(currentScene.ImageWidth)) * float64(math.MaxUint16)),
				B: uint16((float64(y) / float64(currentScene.ImageHeight)) * float64(math.MaxUint16)),
				A: uint16(1.0 * float64(math.MaxUint16))}
			newImage.SetRGBA64(x, y, col)
		}
	}
	imageExists, err := image.DoesExist(plData, sceneID)
	if err != nil {
		fmt.Printf("Error in tracing.go: %s\n", err.Error())
		scene.UpdateRenderStatus(plData, sceneID, renderstatus.Error)
		return
	}
	if imageExists {
		img, err := image.Retrieve(plData, sceneID)
		if err != nil {
			fmt.Printf("Error in tracing.go: %s\n", err.Error())
			scene.UpdateRenderStatus(plData, sceneID, renderstatus.Error)
			return
		}
		img.ImageData = newImage
		image.Update(plData, img)
	} else {
		img := &image.Image{
			SceneID:   sceneID,
			ImageData: newImage,
		}
		_, err := image.Create(plData, img)
		if err != nil {
			fmt.Printf("Error in tracing.go: %s\n", err.Error())
			scene.UpdateRenderStatus(plData, sceneID, renderstatus.Error)
			return
		}
	}
	err = scene.UpdateRenderStatus(plData, sceneID, renderstatus.Completed)
	if err != nil {
		fmt.Printf("Error in tracing.go: %s\n", err.Error())
		scene.UpdateRenderStatus(plData, sceneID, renderstatus.Error)
		return
	}
}
