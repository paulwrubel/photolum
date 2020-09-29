package texture

import "github.com/paulwrubel/photolum/tracing/shading"

// Texture defines behaviors of a Texture implementation
type Texture interface {
	Value(u, v float64) shading.Color
}
