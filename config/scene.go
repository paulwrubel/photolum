package config

// Scene contains informations required to build and save a raytraced scene
type Scene struct {
	ImageWidth  int    `json:"image_width"`  // width of the image in pixels
	ImageHeight int    `json:"image_height"` // height of the image in pixels
	FileType    string `json:"file_type"`    // image file type (png, jpg, etc.)
}
