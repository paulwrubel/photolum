package disk

import (
	"testing"

	"github.com/paulwrubel/photolum/config/geometry"
)

var diskHit bool

func TestDiskIntersectionHit(t *testing.T) {
	disk := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.0,
			Z: 1.0,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	_, h := disk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if !h {
		t.Errorf("Expected true (hit) but got %t\n", h)
	}
}

func BenchmarkDiskIntersectionHit(b *testing.B) {
	disk := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.0,
			Z: 1.0,
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
		_, h = disk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	diskHit = h
}

func TestDiskIntersectionReverseHit(t *testing.T) {
	disk := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: 1.0,
		},
	}
	_, h := disk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if !h {
		t.Errorf("Expected true (hit) but got %t\n", h)
	}
}

func BenchmarkDiskIntersectionReverseHit(b *testing.B) {
	disk := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: 1.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = disk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	diskHit = h
}

func TestDiskIntersectionMiss(t *testing.T) {
	disk := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 1.0,
			Y: 1.0,
			Z: 1.0,
		},
		Direction: geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	_, h := disk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if h {
		t.Errorf("Expected false (miss) but got %t\n", h)
	}
}

func BenchmarkDiskIntersectionMiss(b *testing.B) {
	disk := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 1.0,
			Y: 1.0,
			Z: 1.0,
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
		_, h = disk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	diskHit = h
}

func TestDiskIntersectionParallelMiss(t *testing.T) {
	disk := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.0,
			Z: 1.0,
		},
		Direction: geometry.Vector{
			X: 1.0,
			Y: 0.0,
			Z: 0.0,
		},
	}
	_, h := disk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	if h {
		t.Errorf("Expected false (miss) but got %t\n", h)
	}
}

func BenchmarkDiskIntersectionParallelMiss(b *testing.B) {
	disk := Unit(0.0, 0.0, 0.0)
	r := geometry.Ray{
		Origin: geometry.Point{
			X: 0.0,
			Y: 0.0,
			Z: 1.0,
		},
		Direction: geometry.Vector{
			X: 1.0,
			Y: 0.0,
			Z: 0.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = disk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308, nil)
	}
	diskHit = h
}
