package config

import (
	"github.com/paulwrubel/photolum/config/shading"
	"github.com/paulwrubel/photolum/enumeration/filetype"
)

// Parameters holds top-level information about the program's execution and the image's properties
type Parameters struct {
	ImageWidth               int               // width of the image in pixels
	ImageHeight              int               // height of the image in pixels
	FileType                 filetype.FileType // image file type (png, jpg, etc.)
	GammaCorrection          float64           // how much gamma correction to perform on the image
	UseScalingTruncation     bool              // should the program truncate over-magnitude colors by scaling linearly as opposed to clamping?
	SamplesPerRound          int               // amount of samples to write per rounds
	RoundCount               int               // amount of rounds per render
	TileWidth                int               // width of a tile in pixels
	TileHeight               int               // height of a tile in pixels
	MaxBounces               int               // amount of reflections to check before giving up
	UseBVH                   bool              // should the program generate and use a Bounding Volume Hierarchy?
	BackgroundColorMagnitude float64           // amount to scale bg color by
	BackgroundColor          shading.Color     // color to return when nothing is intersected
	TMin                     float64           // minimum ray "time" to count intersection
	TMax                     float64           // maximum ray "time" to count intersection
	Scene                    *Scene            // Scene reference
}
