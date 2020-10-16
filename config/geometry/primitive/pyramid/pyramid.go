package pyramid

import (
	"fmt"

	"github.com/paulwrubel/photolum/config/geometry"
	"github.com/paulwrubel/photolum/config/geometry/primitive"
	"github.com/paulwrubel/photolum/config/geometry/primitive/aabb"
	"github.com/paulwrubel/photolum/config/geometry/primitive/primitivelist"
	"github.com/paulwrubel/photolum/config/geometry/primitive/rectangle"
	"github.com/paulwrubel/photolum/config/geometry/primitive/triangle"
	"github.com/paulwrubel/photolum/config/shading/material"
)

// Pyramid represents a pyramid geometric shape
type Pyramid struct {
	A      geometry.Point `json:"a"`
	B      geometry.Point `json:"b"`
	Height float64        `json:"height"`
	list   *primitivelist.PrimitiveList
	box    *aabb.AABB
}

// Setup sets up internal fields of a pyramid
func (p *Pyramid) Setup() (*Pyramid, error) {
	if p.Height <= 0 {
		return nil, fmt.Errorf("pyramid height is 0 or negative")
	}
	if p.A.Y != p.B.Y {
		return nil, fmt.Errorf("pyramid is not directed upwards")
	}

	c1 := geometry.MinComponents(p.A, p.B)
	c3 := geometry.MaxComponents(p.A, p.B)
	c2 := geometry.Point{
		X: c1.X,
		Y: c1.Y,
		Z: c3.Z,
	}
	c4 := geometry.Point{
		X: c3.X,
		Y: c1.Y,
		Z: c1.Z,
	}

	base, err := (&rectangle.Rectangle{
		A:                 p.A,
		B:                 p.B,
		IsCulled:          false,
		HasNegativeNormal: true,
	}).Setup()
	if err != nil {
		return nil, err
	}

	diagonalBaseVectorHalf := c1.To(c3).DivScalar(2.0)
	baseCenterPoint := c1.AddVector(diagonalBaseVectorHalf)
	topPoint := baseCenterPoint.AddVector(geometry.VectorUp.MultScalar(p.Height))

	tri1, err := (&triangle.Triangle{
		A:        c1,
		B:        c2,
		C:        topPoint,
		IsCulled: false,
	}).Setup()
	if err != nil {
		return nil, err
	}

	tri2, err := (&triangle.Triangle{
		A:        c2,
		B:        c3,
		C:        topPoint,
		IsCulled: false,
	}).Setup()
	if err != nil {
		return nil, err
	}

	tri3, err := (&triangle.Triangle{
		A:        c3,
		B:        c4,
		C:        topPoint,
		IsCulled: false,
	}).Setup()
	if err != nil {
		return nil, err
	}

	tri4, err := (&triangle.Triangle{
		A:        c4,
		B:        c1,
		C:        topPoint,
		IsCulled: false,
	}).Setup()
	if err != nil {
		return nil, err
	}

	l, err := primitivelist.FromElements(base, tri1, tri2, tri3, tri4)
	if err != nil {
		return nil, err
	}
	b, _ := l.BoundingBox(0, 0)
	return &Pyramid{
		list: l,
		box:  b,
	}, nil
}

// Intersection computer the intersection of this object and a given ray if it exists
func (p *Pyramid) Intersection(ray geometry.Ray, tMin, tMax float64) (*material.RayHit, bool) {
	if p.box.Intersection(ray, tMin, tMax) {
		return p.list.Intersection(ray, tMin, tMax)
	}
	return nil, false
}

// BoundingBox returns an AABB of this object
func (p *Pyramid) BoundingBox(t0, t1 float64) (*aabb.AABB, bool) {
	return p.box, true
}

// SetMaterial sets this object's material
func (p *Pyramid) SetMaterial(m material.Material) {
	p.list.SetMaterial(m)
}

// IsInfinite returns whether this object is infinite
func (p *Pyramid) IsInfinite() bool {
	return false
}

// IsClosed returns whether this object is closed
func (p *Pyramid) IsClosed() bool {
	return true
}

// Copy returns a shallow copy of this object
func (p *Pyramid) Copy() primitive.Primitive {
	newP := *p
	return &newP
}

// Unit return a unit pyramid
func Unit(xOffset, yOffset, zOffset float64) *Pyramid {
	p, _ := (&Pyramid{
		A: geometry.Point{
			X: 0.0 + xOffset,
			Y: 0.0 + yOffset,
			Z: 0.0 + zOffset,
		},
		B: geometry.Point{
			X: 1.0 + xOffset,
			Y: 0.0 + yOffset,
			Z: 1.0 + zOffset,
		},
		Height: 1.0,
	}).Setup()
	return p
}
