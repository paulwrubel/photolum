package material

import (
	"github.com/paulwrubel/photolum/config/geometry"
	"github.com/paulwrubel/photolum/config/shading"
	"github.com/paulwrubel/photolum/config/shading/texture"
)

// Isotropic represents a volumetric material
// (scatters uniformly in unit sphere)
type Isotropic struct {
	ReflectanceTexture texture.Texture `json:"-"`
	EmittanceTexture   texture.Texture `json:"-"`
}

// Reflectance returns the reflective color at texture coordinates (u, v)
func (i Isotropic) Reflectance(u, v float64) shading.Color {
	return i.ReflectanceTexture.Value(u, v)
}

// Emittance returns the emissive color at texture coordinates (u, v)
func (i Isotropic) Emittance(u, v float64) shading.Color {
	return i.EmittanceTexture.Value(u, v)
}

// IsSpecular returns whether this material is specular in nature (vs. diffuse)
// This is currently unused and is likely to be deprecated in the future
func (i Isotropic) IsSpecular() bool {
	return false
}

// Scatter returns an incoming ray given a RayHit representing the outgoing ray
func (i Isotropic) Scatter(rayHit RayHit) (geometry.Ray, bool) {
	hitPoint := rayHit.Ray.PointAt(rayHit.Time)
	direction := geometry.RandomInUnitSphere()
	return geometry.Ray{
		Origin:    hitPoint,
		Direction: direction,
	}, true
}
