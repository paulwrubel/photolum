package config

import (
	"image"

	"github.com/paulwrubel/photolum/enumeration/filetype"
)

type TracingPayload struct {
	FileType filetype.FileType
	Image    image.Image
}
