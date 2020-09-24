package persistence

import (
	"fmt"
	"image"
	"sync"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
)

var configs = sync.Map{}
var images = sync.Map{}

func SaveConfig(scene config.Scene) (uuid.UUID, error) {
	sceneID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}
	configs.Store(sceneID, scene)
	return sceneID, nil
}

func GetConfig(sceneID uuid.UUID) (config.Scene, error) {
	sceneConfigRaw, ok := configs.Load(sceneID)
	sceneConfig := sceneConfigRaw.(config.Scene)
	if !ok {
		return config.Scene{}, fmt.Errorf("no scene configs founds with sceneID = %s", sceneID.String())
	}
	return sceneConfig, nil
}

func UpdateConfig(sceneID uuid.UUID, scene config.Scene) {
	configs.Store(sceneID, scene)
}

func DeleteConfig(sceneID uuid.UUID) {
	configs.Delete(sceneID)
}

func GetImage(sceneID uuid.UUID) (*image.RGBA64, error) {
	sceneImageRaw, ok := images.Load(sceneID)
	sceneImage := sceneImageRaw.(*image.RGBA64)
	if !ok {
		return nil, fmt.Errorf("no scene images found with sceneID = %s", sceneID.String())
	}
	return sceneImage, nil
}

func UpdateImage(sceneID uuid.UUID, image *image.RGBA64) {
	images.Store(sceneID, image)
}

func DeleteImage(sceneID uuid.UUID) {
	images.Delete(sceneID)
}
