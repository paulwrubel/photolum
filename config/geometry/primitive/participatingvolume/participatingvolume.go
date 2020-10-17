package participatingvolume

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/paulwrubel/photolum/config/geometry"
	"github.com/paulwrubel/photolum/config/geometry/primitive"
	"github.com/paulwrubel/photolum/config/geometry/primitive/aabb"
	"github.com/paulwrubel/photolum/config/shading/material"
)

// ParticipatingVolume represents a participating volume geometry primitive (smoke, fog, fire, etc.)
type ParticipatingVolume struct {
	Density                float64
	Primitive              primitive.Primitive
	negativeInverseDensity float64
	mat                    material.Material
}

// Setup sets up a participating volume
func (pv *ParticipatingVolume) Setup() (*ParticipatingVolume, error) {
	if pv.Density <= 0.0 {
		return nil, fmt.Errorf("Density must be greater than zero")
	}
	pv.negativeInverseDensity = -1.0 / pv.Density
	return pv, nil
}

// Intersection computer the intersection of this primitive and a given ray
func (pv *ParticipatingVolume) Intersection(ray geometry.Ray, tMin, tMax float64, rng *rand.Rand) (*material.RayHit, bool) {

	// hit first part of surface
	rayHit1, wasHit := pv.Primitive.Intersection(ray, -math.MaxFloat64, math.MaxFloat64, rng)
	if !wasHit {
		return nil, false
	}
	// hit second part of surface
	rayHit2, wasHit := pv.Primitive.Intersection(ray, rayHit1.Time+0.0001, math.MaxFloat64, rng)
	if !wasHit {
		return nil, false
	}

	if rayHit1.Time < tMin {
		rayHit1.Time = tMin
	}
	if rayHit2.Time > tMax {
		rayHit2.Time = tMax
	}

	if rayHit1.Time >= rayHit2.Time {
		return nil, false
	}

	if rayHit1.Time < 0 {
		rayHit1.Time = 0
	}

	rayMagnitude := ray.Direction.Magnitude()
	distanceInsideBoundary := (rayHit2.Time - rayHit1.Time) * rayMagnitude
	hitDistance := pv.negativeInverseDensity * math.Log(rng.Float64())

	if hitDistance > distanceInsideBoundary {
		return nil, false
	}

	time := rayHit1.Time + (hitDistance / rayMagnitude)

	return &material.RayHit{
		Ray:         ray,
		NormalAtHit: geometry.VectorRight, // arbitrary, because material is always isotropic
		Time:        time,
		U:           0.0,
		V:           0.0,
		Material:    pv.mat,
	}, true
}

// BoundingBox returns an AABB of this object
func (pv *ParticipatingVolume) BoundingBox(t0, t1 float64) (*aabb.AABB, bool) {
	return pv.Primitive.BoundingBox(t0, t1)
}

// SetMaterial sets this object's material
func (pv *ParticipatingVolume) SetMaterial(m material.Material) {
	pv.mat = m
}

// IsInfinite returns whether this object is infinite
func (pv *ParticipatingVolume) IsInfinite() bool {
	return pv.Primitive.IsInfinite()
}

// IsClosed returns whether this object is closed
func (pv *ParticipatingVolume) IsClosed() bool {
	return pv.Primitive.IsClosed()
}

// Copy returns a shallow copy of this object
func (pv *ParticipatingVolume) Copy() primitive.Primitive {
	newPV := *pv
	return &newPV
}
