package scenes

import "github.com/paulwrubel/photolum/config"

// SceneIDGetRequest contains the sceneID GET endpoint request
type SceneIDGetResponse struct {
	Scene config.Scene `json:"scene"`
}
