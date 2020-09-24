package scenes

import "github.com/paulwrubel/photolum/config"

// ScenePostRequest contains the scene POST endpoint request
type ScenePostRequest struct {
	Scene config.Scene `json:"scene"`
}
