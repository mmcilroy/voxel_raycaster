package voxel

import (
	"fmt"
)

type VoxelGrid struct {
	Parent     *VoxelGrid
	NumVoxelsX int
	NumVoxelsY int
	NumVoxelsZ int
	VoxelSize  float32
	Voxels     []uint64
}

func NewVoxelGrid(nx, ny, nz int, sz float32) *VoxelGrid {
	nv := nx / 4 * ny / 4 * nz / 4
	if nv == 0 {
		nv = 1
	}
	return &VoxelGrid{
		NumVoxelsX: nx,
		NumVoxelsY: ny,
		NumVoxelsZ: nz,
		VoxelSize:  sz,
		Voxels:     make([]uint64, nv),
	}
}

func voxelBitMask(x, y, z int) uint64 {
	bit := uint64(x + (z * 4) + (y * 4 * 4))
	mask := uint64(1) << bit
	return mask
}

func (grid *VoxelGrid) VoxelIndex(x, y, z int) int {
	vx, vy, vz := x/4, y/4, z/4
	return vx + vz*grid.NumVoxelsX/4 + vy*grid.NumVoxelsX/4*grid.NumVoxelsZ/4
}

func (grid *VoxelGrid) GetVoxel(x, y, z int) bool {
	i := grid.VoxelIndex(x, y, z)
	if i < 0 || i >= len(grid.Voxels) {
		fmt.Printf("GetVoxel: Invalid XYZ - %d, %d, %d\n", x, y, z)
		return false
	}

	voxel := grid.Voxels[i]

	if voxel == 0 {
		return false
	}

	mask := voxelBitMask(x%4, y%4, z%4)
	return voxel&mask != 0
}

func (grid *VoxelGrid) SetVoxel(x, y, z int, set bool) {
	i := grid.VoxelIndex(x, y, z)
	if i < 0 || i >= len(grid.Voxels) {
		fmt.Printf("SetVoxel: Invalid XYZ - %d, %d, %d\n", x, y, z)
		return
	}

	voxel := &grid.Voxels[i]
	mask := voxelBitMask(x%4, y%4, z%4)

	if set {
		*voxel = *voxel | mask
	} else {
		*voxel = *voxel & ^mask
	}
}

func (grid *VoxelGrid) Clear() {
	for i := 0; i < len(grid.Voxels); i++ {
		grid.Voxels[i] = 0
	}
}

func (grid *VoxelGrid) Compress() *VoxelGrid {
	// create the new grid which will be half the size
	newGrid := NewVoxelGrid(grid.NumVoxelsX/2, grid.NumVoxelsY/2, grid.NumVoxelsZ/2, grid.VoxelSize*2)
	newGrid.Parent = grid

	for y := 0; y < newGrid.NumVoxelsY; y++ {
		for z := 0; z < newGrid.NumVoxelsZ; z++ {
			for x := 0; x < newGrid.NumVoxelsX; x++ {

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
