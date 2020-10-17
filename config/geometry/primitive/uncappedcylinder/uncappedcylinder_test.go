package uncappedcylinder

import (
	"testing"

	"github.com/paulwrubel/photolum/config/geometry"
)

var ucHit bool

func TestUncappedCylinderIntersectionHit(t *testing.T) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.5,
			Z: 1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	_, h := uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if !h {
		t.Errorf("Expected true (hit) but got %t\n", h)
	}
}

func BenchmarkUncappedCylinderIntersectionHit(b *testing.B) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.5,
			Z: 1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	ucHit = h
}

func TestUncappedCylinderIntersectionSecondHit(t *testing.T) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.5,
			Z: 0.0,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	_, h := uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if !h {
		t.Errorf("Expected true (hit) but got %t\n", h)
	}
}

func BenchmarkUncappedCylinderIntersectionSecondHit(b *testing.B) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.5,
			Z: 0.0,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	ucHit = h
}

func TestUncappedCylinderIntersectionSideMiss(t *testing.T) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 1.5,
			Y: 0.5,
			Z: 1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	_, h := uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if h {
		t.Errorf("Expected false (miss) but got %t\n", h)
	}
}

func BenchmarkUncappedCylinderIntersectionSideMiss(b *testing.B) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 1.5,
			Y: 0.5,
			Z: 1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	ucHit = h
}

func TestUncappedCylinderIntersectionBehindMiss(t *testing.T) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.5,
			Z: -1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	_, h := uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if h {
		t.Errorf("Expected false (miss) but got %t\n", h)
	}
}

func BenchmarkUncappedCylinderIntersectionBehindMiss(b *testing.B) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.5,
			Z: -1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	ucHit = h
}

func TestUncappedCylinderIntersectionTopMiss(t *testing.T) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 1.5,
			Z: 1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	_, h := uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if h {
		t.Errorf("Expected false (miss) but got %t\n", h)
	}
}

func BenchmarkUncappedCylinderIntersectionTopMiss(b *testing.B) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 1.5,
			Z: 1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	ucHit = h
}

func TestUncappedCylinderIntersectionBottomMiss(t *testing.T) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: -0.5,
			Z: 1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	_, h := uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if h {
		t.Errorf("Expected false (miss) but got %t\n", h)
	}
}

func BenchmarkUncappedCylinderIntersectionBottomMiss(b *testing.B) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: -0.5,
			Z: 1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	ucHit = h
}

func TestUncappedCylinderIntersectionParallelMiss(t *testing.T) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.5,
			Z: 1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 1.0,
			Z: 0.0,
		},
	}
	_, h := uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if h {
		t.Errorf("Expected false (miss) but got %t\n", h)
	}
}

func BenchmarkUncappedCylinderIntersectionParallelMiss(b *testing.B) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.5,
			Z: 1.5,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 1.0,
			Z: 0.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	ucHit = h
}

func TestUncappedCylinderIntersectionInsideParallelMiss(t *testing.T) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.5,
			Z: 0.0,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 1.0,
			Z: 0.0,
		},
	}
	_, h := uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if h {
		t.Errorf("Expected false (miss) but got %t\n", h)
	}
}

func BenchmarkUncappedCylinderIntersectionInsideParallelMiss(b *testing.B) {
	uc := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.5,
			Z: 0.0,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 1.0,
			Z: 0.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = uc.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	ucHit = h
}
