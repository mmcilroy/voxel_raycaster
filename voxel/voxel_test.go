package voxel

import (
	"testing"
)

func TestVoxelGetSet(t *testing.T) {
	grid := NewVoxelGrid(4, 4, 4, 1)

	for y := int32(0); y < grid.NumVoxelsY-1; y++ {
		for z := int32(0); z < grid.NumVoxelsZ-1; z++ {
			for x := int32(0); x < grid.NumVoxelsX-1; x++ {

				prev := grid.Voxels[grid.VoxelIndex(x, y, z)]
				grid.SetVoxel(x, y, z, true)
				curr := grid.Voxels[grid.VoxelIndex(x, y, z)]

				// every voxel set should increase the overall value
				if curr <= prev {
					t.FailNow()
				}
			}
		}
	}

	for y := int32(0); y < grid.NumVoxelsY-1; y++ {
		for z := int32(0); z < grid.NumVoxelsZ-1; z++ {
			for x := int32(0); x < grid.NumVoxelsX-1; x++ {

				prev := grid.Voxels[grid.VoxelIndex(x, y, z)]
				grid.SetVoxel(x, y, z, false)
				curr := grid.Voxels[grid.VoxelIndex(x, y, z)]

				// every voxel cleared should decrease the overall value
				if curr >= prev {
					t.FailNow()
				}
			}
		}
	}
}

/*
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
*/
