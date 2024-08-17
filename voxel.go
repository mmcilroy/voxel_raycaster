package main

import "fmt"

type VoxelGrid struct {
	Size   int
	Array  []uint64
	Scale  float32
	Parent *VoxelGrid
}

func voxelBitMask(x, y, z int) uint64 {
	bit := uint64(x + (z * 4) + (y * 4 * 4))
	mask := uint64(1) << bit
	return mask
}

func voxelSet(voxel *uint64, x, y, z int, set bool) {
	if x > 4 || y > 4 || z > 4 {
		fmt.Printf("voxelSet: Invalid XYZ - %d,%d,%d\n", x, y, z)
		return
	}
	mask := voxelBitMask(x, y, z)
	if set {
		*voxel = *voxel | mask
	} else {
		*voxel = *voxel & ^mask
	}
}

func voxelGet(voxel uint64, x, y, z int) bool {
	if x > 4 || y > 4 || z > 4 {
		fmt.Printf("voxelGet: Invalid XYZ - %d,%d,%d\n", x, y, z)
		return false
	}
	mask := voxelBitMask(x, y, z)
	return voxel&mask != 0
}

// create a new voxel grid
// size - number of voxels along one axis
// for example size 4 would result in 4x4x4 voxels (64)
func NewVoxelGrid(size int) *VoxelGrid {
	return &VoxelGrid{
		Size:  size,
		Array: make([]uint64, size/4*size/4*size/4),
		Scale: 1,
	}
}

func (grid *VoxelGrid) Clear() {
	for i := 0; i < len(grid.Array); i++ {
		grid.Array[i] = 0
	}
}

// set the value (true/false) of a voxel with the VoxelGrid
// x, y, z - must be with range (0 .. size - 1)
// v - true or false
func (grid *VoxelGrid) SetVoxel(x, y, z int, v bool) {
	i := grid.GetVoxelIndex(x, y, z)
	if i < 0 || i >= len(grid.Array) {
		fmt.Printf("SetVoxel: Invalid XYZ - %d, %d, %d\n", x, y, z)
		return
	}
	voxelSet(&grid.Array[i], x%4, y%4, z%4, v)
}

// pretty obvious
func (grid *VoxelGrid) GetVoxel(x, y, z int) bool {
	i := grid.GetVoxelIndex(x, y, z)
	if i < 0 || i >= len(grid.Array) {
		fmt.Printf("GetVoxel: Invalid XYZ - %d, %d, %d\n", x, y, z)
		return false
	}
	v := grid.Array[i]
	if v == 0 {
		return false
	}
	return voxelGet(v, x%4, y%4, z%4)
}

func (grid *VoxelGrid) Compress() *VoxelGrid {
	// create the new grid which will be 4 times smaller
	newGrid := NewVoxelGrid(grid.Size / 4)
	newGrid.Parent = grid
	newGrid.Scale = grid.Scale / 4

	// iterate over the current grid
	for y := 0; y < grid.Size; y++ {
		for z := 0; z < grid.Size; z++ {
			for x := 0; x < grid.Size; x++ {
				// directly access the array element at x, y, z, which corresponds to a 4x4x4 grid
				// if the value is not zero then set the corresponding voxel in our new grid
				i := grid.GetVoxelIndex(x, y, z)
				newGrid.SetVoxel(x/4, y/4, z/4, grid.Array[i] > 0)
			}
		}
	}

	return newGrid
}

func (grid *VoxelGrid) Compress2() *VoxelGrid {
	// create the new grid which will be half the size
	newGrid := NewVoxelGrid(grid.Size / 2)
	newGrid.Parent = grid
	newGrid.Scale = grid.Scale / 2

	for y := 0; y < newGrid.Size; y++ {
		for z := 0; z < newGrid.Size; z++ {
			for x := 0; x < newGrid.Size; x++ {

				px, py, pz := x*2, y*2, z*2

				s := grid.GetVoxel(px, py, pz) ||
					grid.GetVoxel(px+1, py, pz) ||
					grid.GetVoxel(px, py+1, pz) ||
					grid.GetVoxel(px+1, py+1, pz) ||
					grid.GetVoxel(px, py, pz+1) ||
					grid.GetVoxel(px+1, py, pz+1) ||
					grid.GetVoxel(px, py+1, pz+1) ||
					grid.GetVoxel(px+1, py+1, pz+1)

				newGrid.SetVoxel(x, y, z, s)
			}
		}
	}

	return newGrid
}

func (grid *VoxelGrid) GetVoxelIndex(x, y, z int) int {
	vx, vy, vz := x/4, y/4, z/4
	return vx + vz*grid.Size/4 + vy*grid.Size/4*grid.Size/4
}
