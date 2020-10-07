package box

import (
	"fmt"

	"github.com/paulwrubel/photolum/config/geometry"
	"github.com/paulwrubel/photolum/config/geometry/primitive"
	"github.com/paulwrubel/photolum/config/geometry/primitive/aabb"
	"github.com/paulwrubel/photolum/config/geometry/primitive/primitivelist"
	"github.com/paulwrubel/photolum/config/geometry/primitive/rectangle"
	"github.com/paulwrubel/photolum/config/shading/material"
)

// Box represents a box
type Box struct {
	A                  geometry.Point `json:"a"`
	B                  geometry.Point `json:"b"`
	HasInvertedNormals bool           `json:"has_inverted_normals"`
	list               *primitivelist.PrimitiveList
	box                *aabb.AABB
}

// type Data struct {
// }

// Setup sets up this object's internal field
func (b *Box) Setup() (*Box, error) {
	c1 := geometry.MinComponents(b.A, b.B)
	c8 := geometry.MaxComponents(b.A, b.B)

	if c1.X == c8.X || c1.Y == c8.Y || c1.Z == c8.Z {
		return nil, fmt.Errorf("box resolves to point, line, or plane")
	}

	rNegX, err := (&rectangle.Rectangle{
		A: c1,
		B: geometry.Point{
			X: c1.X,
			Y: c8.Y,
			Z: c8.Z,
		},
		HasNegativeNormal: !b.HasInvertedNormals,
	}).Setup()
	if err != nil {
		return nil, err
	}

	rPosX, err := (&rectangle.Rectangle{
		A: geometry.Point{
			X: c8.X,
			Y: c1.Y,
			Z: c1.Z,
		},
		B:                 c8,
		HasNegativeNormal: b.HasInvertedNormals,
	}).Setup()
	if err != nil {
		return nil, err
	}

	rNegY, err := (&rectangle.Rectangle{
		A: c1,
		B: geometry.Point{
			X: c8.X,
			Y: c1.Y,
			Z: c8.Z,
		},
		HasNegativeNormal: !b.HasInvertedNormals,
	}).Setup()
	if err != nil {
		return nil, err
	}

	rPosY, err := (&rectangle.Rectangle{
		A: geometry.Point{
			X: c1.X,
			Y: c8.Y,
			Z: c1.Z,
		},
		B:                 c8,
		HasNegativeNormal: b.HasInvertedNormals,
	}).Setup()
	if err != nil {
		return nil, err
	}

	rNegZ, err := (&rectangle.Rectangle{
		A: c1,
		B: geometry.Point{
			X: c8.X,
			Y: c8.Y,
			Z: c1.Z,
		},
		HasNegativeNormal: !b.HasInvertedNormals,
	}).Setup()
	if err != nil {
		return nil, err
	}

	rPosZ, err := (&rectangle.Rectangle{
		A: geometry.Point{
			X: c1.X,
			Y: c1.Y,
			Z: c8.Z,
		},
		B:                 c8,
		HasNegativeNormal: b.HasInvertedNormals,
	}).Setup()
	if err != nil {
		return nil, err
	}

	primitiveList, err := primitivelist.FromElements(rNegX, rPosX, rNegY, rPosY, rNegZ, rPosZ)
	if err != nil {
		return nil, err
	}

	b.list = primitiveList
	b.box, _ = primitiveList.BoundingBox(0, 0)
	return b, nil
}

// Intersection computer the intersection of this object and a given ray if it exists
func (b *Box) Intersection(ray geometry.Ray, tMin, tMax float64) (*material.RayHit, bool) {
	if b.box.Intersection(ray, tMin, tMax) {
		return b.list.Intersection(ray, tMin, tMax)
	}
	return nil, false
}

// BoundingBox returns an AABB for this object
func (b *Box) BoundingBox(t0, t1 float64) (*aabb.AABB, bool) {
	return b.box, true
}

// SetMaterial sets this object's material
func (b *Box) SetMaterial(m material.Material) {
	b.list.SetMaterial(m)
}

// IsInfinite returns whether this object is infinite
func (b *Box) IsInfinite() bool {
	return false
}

// IsClosed returns whether this object is closed
func (b *Box) IsClosed() bool {
	return true
}

// Copy returns a shallow copy of this object
func (b *Box) Copy() primitive.Primitive {
	newB := *b
	return &newB
}

// Unit returns a unit box
func Unit(xOffset, yOffset, zOffset float64) *Box {
	b, _ := (&Box{
		A: geometry.Point{
			X: 0.0 + xOffset,
			Y: 0.0 + yOffset,
			Z: 0.0 + zOffset,
		},
		B: geometry.Point{
			X: 1.0 + xOffset,
			Y: 1.0 + yOffset,
			Z: 1.0 + zOffset,
		},
	}).Setup()
	return b
}
