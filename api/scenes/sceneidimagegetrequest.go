package scenes

// SceneIDGetRequest contains the sceneID GET endpoint request
type SceneIDImageGetRequest struct {
	SceneID  string `json:"scene_id"`
	Protocol string `json:"protocol"`
}
