package material

import (
	"math/rand"

	"github.com/paulwrubel/photolum/tracing/geometry"
	"github.com/paulwrubel/photolum/tracing/shading"
)

// Material described the implementation of a surface material
type Material interface {
	Reflectance(u, v float64) shading.Color
	Emittance(u, v float64) shading.Color
	IsSpecular() bool
	Scatter(RayHit, *rand.Rand) (geometry.Ray, bool)
}

// RayHit is a loose gathering of information about a ray's intersection with a surface
type RayHit struct {
	Ray         geometry.Ray
	NormalAtHit geometry.Vector
	Time        float64
	U           float64 // texture coordinate U
	V           float64 // texture coordinate V
	Material    Material
}
