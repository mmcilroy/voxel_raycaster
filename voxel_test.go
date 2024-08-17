package main

import (
	"testing"
)

func TestVoxelBitmasks(t *testing.T) {
	voxel := uint64(0)

	for y := 0; y < 4; y++ {
		for z := 0; z < 4; z++ {
			for x := 0; x < 4; x++ {
				prev := voxel
				voxelSet(&voxel, x, y, z, true)

				// every voxel set should increase the overall value
				if !(voxel > prev) {
					t.FailNow()
				}
			}
		}
	}

	for y := 0; y < 4; y++ {
		for z := 0; z < 4; z++ {
			for x := 0; x < 4; x++ {
				prev := voxel
				voxelSet(&voxel, x, y, z, false)

				// every voxel cleared should decrease the overall value
				if !(voxel < prev) {
					t.FailNow()
				}
			}
		}
	}
}

func TestVoxelBasic(t *testing.T) {
	voxels := NewVoxelGrid(4)

	if voxels.Size != 4 {
		t.FailNow()
	}

	if len(voxels.Array) != 1 {
		t.FailNow()
	}

	voxels = NewVoxelGrid(16)

	if voxels.Size != 16 {
		t.FailNow()
	}

	if len(voxels.Array) != 4*4*4 {
		t.FailNow()
	}
}

func TestVoxelCompression(t *testing.T) {
	grid := NewVoxelGrid(16)

	for y := 0; y < grid.Size; y++ {
		for z := 0; z < grid.Size; z++ {
			for x := 0; x < grid.Size; x++ {
				if y < grid.Size/2 {
					grid.SetVoxel(x, y, z, true)
				}
			}
		}
	}

	compressed := grid.Compress()

	if compressed.Size != 4 {
		t.FailNow()
	}

	for y := 0; y < compressed.Size; y++ {
		for z := 0; z < compressed.Size; z++ {
			for x := 0; x < compressed.Size; x++ {
				voxel := compressed.GetVoxel(x, y, z)
				if y < 2 && voxel == false {
					t.FailNow()
				}
				if y > 2 && voxel == true {
					t.FailNow()
				}
			}
		}
	}
}
