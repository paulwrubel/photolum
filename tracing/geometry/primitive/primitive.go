package primitive

import (
	"github.com/paulwrubel/photolum/tracing/geometry"
	"github.com/paulwrubel/photolum/tracing/geometry/primitive/aabb"
	"github.com/paulwrubel/photolum/tracing/shading/material"
)

// Primitive represents a geometry object with a material in 3D space in the scene
type Primitive interface {
	Intersection(geometry.Ray, float64, float64) (*material.RayHit, bool)
	BoundingBox(float64, float64) (*aabb.AABB, bool)
	SetMaterial(material.Material)
	IsInfinite() bool
	IsClosed() bool
	Copy() Primitive
}
