package material

import (
	"github.com/paulwrubel/photolum/config/geometry"
	"github.com/paulwrubel/photolum/config/shading"
	"github.com/paulwrubel/photolum/config/shading/texture"
)

// Metal is an implementation of a Material
// It represents a perfect or near-perfect specularly reflective material
type Metal struct {
	ReflectanceTexture texture.Texture `json:"-"`
	EmittanceTexture   texture.Texture `json:"-"`
	Fuzziness          float64         `json:"fuzziness"`
}

// Reflectance returns the reflective color at texture coordinates (u, v)
func (m Metal) Reflectance(u, v float64) shading.Color {
	return m.ReflectanceTexture.Value(u, v)
}

// Emittance returns the emissive color at texture coordinates (u, v)
func (m Metal) Emittance(u, v float64) shading.Color {
	return m.EmittanceTexture.Value(u, v)
}

// IsSpecular returns whether this material is specular in nature (vs. diffuse)
// This is currently unused and is likely to be deprecated in the future
func (m Metal) IsSpecular() bool {
	return true
}

// Scatter returns an incoming ray given a RayHit representing the outgoing ray
func (m Metal) Scatter(rayHit RayHit) (geometry.Ray, bool) {
	hitPoint := rayHit.Ray.PointAt(rayHit.Time)
	normal := rayHit.NormalAtHit

	reflectionVector := rayHit.Ray.Direction.Unit().ReflectAround(normal)
	reflectionVector = reflectionVector.Add(geometry.RandomInUnitSphere().MultScalar(m.Fuzziness))
	if reflectionVector.Dot(normal) > 0 {
		return geometry.Ray{
			Origin:    hitPoint,
			Direction: reflectionVector,
		}, true
	}
	return geometry.RayZero, false
}
