package texture

import (
	"bytes"
	"image"

	"github.com/paulwrubel/photolum/config/shading"
)

// Image holds information about a texture based on an image
type Image struct {
	ImageData []byte      `json:"-"`
	Gamma     float64     `json:"gamma"`
	Magnitude float64     `json:"magnitude"`
	Image     image.Image `json:"-"`
}

// Load decodes the image from the given filename and performs other setup actions
func (it *Image) Load() error {
	var err error
	it.Image, _, err = image.Decode(bytes.NewReader(it.ImageData))
	if err != nil {
		return err
	}
	return nil
}

// Value returns the color of the image at the given texture coordinates
// parameters u and v have a valid range [0.0, 1.0)
func (it *Image) Value(u, v float64) shading.Color {
	// convert to image coordinates
	x := int(u * float64(it.Image.Bounds().Dx()-1))
	y := int((1.0 - v) * float64(it.Image.Bounds().Dy()-1))
	// get the color of the image at that point
	color := it.Image.At(x, y)
	// convert to a color, de-gamma, and apply magnitude
	return shading.MakeColor(color).Pow(it.Gamma).MultScalar(it.Magnitude)
}
