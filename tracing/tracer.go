package tracing

import (
	"image"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/persistence"
)

func SaveImage(sceneID uuid.UUID) error {
	scene, err := persistence.GetConfig(sceneID)
	if err != nil {
		return err
	}
	newImage := image.NewRGBA64(image.Rect(0, 0, scene.ImageWidth, scene.ImageHeight))
	persistence.UpdateImage(sceneID, newImage)
	return nil
}
