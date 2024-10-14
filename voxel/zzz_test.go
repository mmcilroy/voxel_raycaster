package voxel

import (
	"fmt"
	"testing"
)

func TestTraceHit(t *testing.T) {
	voxels := NewTestVoxels(16, 16, 16, 1)
	voxels.Set(8, 8, 8, true)

	for y := 0; y < 16; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {

				if x == 8 && y == 8 && z == 8 {
					continue
				}

				rs := Vector3f{X: 0.5 + float32(x), Y: 0.5 + float32(y), Z: 0.5 + float32(z)}
				re := Vector3f{X: 8.5, Y: 8.5, Z: 8.5}
				rd := Direction(re, rs)

				params := TraceParams{
					RayStart: rs,
					RayDir:   rd,
					MaxSteps: 0,
				}

				result := Trace(&voxels, params)

				if !result.Hit {
					t.Fatalf("Failed to hit: %d %d %d %+v\n", x, y, z, result)
				}

				if result.OOB {
					t.Fatalf("Incorrect oob: %d %d %d %+v\n", x, y, z, result)
				}

				if !result.MapPos.Equals(Vector3i{X: 8, Y: 8, Z: 8}) {
					t.Fatalf("Incorrect mapPos: %d %d %d %+v\n", x, y, z, result)
				}

				if result.NumSteps < 1 {
					t.Fatalf("Suspicious numSteps: %d %d %d %+v\n", x, y, z, result)
				}

				distance := Distance(result.HitPos, re)
				if distance > 0.867 {
					t.Fatalf("Suspicious distance: %d %d %d %f %+v\n", x, y, z, distance, result)
				}
			}
		}
	}
}

func TestTraceMiss(t *testing.T) {
	voxels := NewTestVoxels(16, 16, 16, 1)

	rs := Vector3f{X: 0.5, Y: 0.5, Z: 0.5}
	re := Vector3f{X: 8.5, Y: 8.5, Z: 8.5}
	rd := Direction(re, rs)

	params := TraceParams{
		RayStart: rs,
		RayDir:   rd,
		MaxSteps: 0,
		Callback: func(voxels *Voxels, mapPos Vector3i) {
			fmt.Printf("Check: %d %d %d %.02f\n", mapPos.X, mapPos.Y, mapPos.Z, (*voxels).Size())
		},
	}

	result := Trace(&voxels, params)

	if result.Hit {
		t.Fatalf("Incorrect hit: %+v\n", result)
	}

	if !result.OOB {
		t.Fatalf("Incorrect oob: %+v\n", result)
	}

	if result.NumSteps < 1 {
		t.Fatalf("Suspicious numSteps: %+v\n", result)
	}
}

func TestMipmap(t *testing.T) {

	voxels0 := NewTestVoxels(64, 64, 64, 1)
	voxels0.Set(32, 32, 32, true)

	voxels1 := voxels0.Compress()
	voxels2 := (*voxels1).Compress()
	voxels3 := (*voxels2).Compress()

	rs := Vector3f{X: 0.5, Y: 0.5, Z: 0.5}
	re := Vector3f{X: 32.5, Y: 34.5, Z: 32.5}
	rd := Direction(re, rs)

	tracer := MipmapTracerImpl{
		Voxels: []*Voxels{&voxels0, voxels1, voxels2, voxels3},
	}

	params := TraceParams{
		RayStart: rs,
		RayDir:   rd,
		MaxSteps: 0,
		Callback: func(voxels *Voxels, mapPos Vector3i) {
			fmt.Printf("Check: %d %d %d %.02f\n", mapPos.X, mapPos.Y, mapPos.Z, (*voxels).Size())
		},
	}

	fmt.Printf("%+v\n", tracer.Trace(params))
}
