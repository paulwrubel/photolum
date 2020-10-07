package texture

import "github.com/paulwrubel/photolum/config/shading"

// Color holds information about a solid-colored texture
type Color struct {
	Color shading.Color `json:"color"`
}

// Value returns a color at a given texture coordinate
// this value is always the same, as the color is solid
func (ct *Color) Value(u, v float64) shading.Color {
	return ct.Color
}
