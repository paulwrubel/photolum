package tracing

import (
	"image"
	"image/color"
	"math"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config/renderstatus"
	"github.com/paulwrubel/photolum/persistence"
)

func StartRender(sceneID uuid.UUID) {
	go TraceImage(sceneID)
	persistence.UpdateRenderStatus(sceneID, renderstatus.Running)
}

func StopRender(sceneID uuid.UUID) {
	persistence.UpdateRenderStatus(sceneID, renderstatus.Stopping)
}

func TraceImage(sceneID uuid.UUID) {
	sceneData, err := persistence.Retrieve(sceneID)
	if err != nil {
		persistence.UpdateRenderStatus(sceneID, renderstatus.Error)
		return
	}
	scene := sceneData.Scene
	newImage := image.NewRGBA64(image.Rect(0, 0, scene.ImageWidth, scene.ImageHeight))
	for y := 0; y < scene.ImageHeight; y++ {
		for x := 0; x < scene.ImageWidth; x++ {
			col := color.RGBA64{
				R: uint16(0.0 * float64(math.MaxUint16)),
				G: uint16((float64(x) / float64(scene.ImageWidth)) * float64(math.MaxUint16)),
				B: uint16((float64(y) / float64(scene.ImageHeight)) * float64(math.MaxUint16)),
				A: uint16(1.0 * float64(math.MaxUint16))}
			newImage.SetRGBA64(x, y, col)
		}
	}
	sceneData.Image = newImage
	persistence.Update(sceneID, sceneData)
	persistence.UpdateRenderStatus(sceneID, renderstatus.Completed)
}
