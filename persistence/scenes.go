package persistence

import (
	"fmt"
	"image"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
)

type SceneData struct {
	SceneID      uuid.UUID
	createdTime  time.Time
	modifiedTime time.Time
	accessedTime time.Time
	Scene        config.Scene
	Image        *image.RGBA64
}

var data = sync.Map{}

func Create(sceneData SceneData) (uuid.UUID, error) {
	if sceneData.SceneID == uuid.Nil {
		newSceneID, err := uuid.NewRandom()
		if err != nil {
			return uuid.Nil, err
		}
		sceneData.SceneID = newSceneID
	}
	timeNow := time.Now()
	sceneData.createdTime = timeNow
	sceneData.modifiedTime = timeNow
	sceneData.accessedTime = timeNow
	data.Store(sceneData.SceneID, sceneData)
	return sceneData.SceneID, nil
}

func Retrieve(sceneID uuid.UUID) (SceneData, error) {
	sceneDataRaw, ok := data.Load(sceneID)
	if !ok {
		return SceneData{}, fmt.Errorf("no scene data found with sceneID = %s", sceneID.String())
	}
	sceneData := sceneDataRaw.(SceneData)
	sceneData.accessedTime = time.Now()
	data.Store(sceneData.SceneID, sceneData)
	return sceneData, nil
}

func Update(sceneID uuid.UUID, sceneData SceneData) {
	sceneData.modifiedTime = time.Now()
	data.Store(sceneData.SceneID, sceneData)
}

func Delete(sceneID uuid.UUID) {
	data.Delete(sceneID)
}

func RetrieveAll() []SceneData {
	var sceneDataList []SceneData
	data.Range(func(_ interface{}, value interface{}) bool {
		sceneData, _ := value.(SceneData)
		sceneDataList = append(sceneDataList, sceneData)
		return true
	})
	return sceneDataList
}
