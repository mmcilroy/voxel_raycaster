package voxel

import (
	"fmt"
)

type VoxelGrid struct {
	Parent     *VoxelGrid // higher res version
	Child      *VoxelGrid // lower res version
	NumVoxelsX int32
	NumVoxelsY int32
	NumVoxelsZ int32
	VoxelSize  float32
	Voxels     []uint8
}

func NewVoxelGrid(nx, ny, nz int32, sz float32) *VoxelGrid {
	return &VoxelGrid{
		NumVoxelsX: nx,
		NumVoxelsY: ny,
		NumVoxelsZ: nz,
		VoxelSize:  sz,
		Voxels:     make([]uint8, nx/2*ny/2*nz/2),
	}
}

func voxelBitMask(x, y, z int32) uint8 {
	bit := uint8(x + (z * 2) + (y * 2 * 2))
	mask := uint8(1) << bit
	return mask
}

func (grid *VoxelGrid) VoxelIndex(x, y, z int32) int32 {
	vx, vy, vz := x/2, y/2, z/2
	return vx + vz*grid.NumVoxelsX/2 + vy*grid.NumVoxelsX/2*grid.NumVoxelsZ/2
}

func (grid *VoxelGrid) GetVoxel(x, y, z int32) bool {
	i := grid.VoxelIndex(x, y, z)
	if i < 0 || i >= int32(len(grid.Voxels)) {
		fmt.Printf("GetVoxel: Invalid XYZ - %d, %d, %d\n", x, y, z)
		return false
	}

	voxel := grid.Voxels[i]

	if voxel == 0 {
		return false
	}

	mask := voxelBitMask(x%2, y%2, z%2)
	return voxel&mask != 0
}

func (grid *VoxelGrid) SetVoxel(x, y, z int32, set bool) {
	i := grid.VoxelIndex(x, y, z)
	if i < 0 || i >= int32(len(grid.Voxels)) {
		fmt.Printf("SetVoxel: Invalid XYZ - %d, %d, %d\n", x, y, z)
		return
	}

	voxel := &grid.Voxels[i]
	mask := voxelBitMask(x%2, y%2, z%2)

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
	grid.Child = newGrid

	for y := int32(0); y < newGrid.NumVoxelsY; y++ {
		for z := int32(0); z < newGrid.NumVoxelsZ; z++ {
			for x := int32(0); x < newGrid.NumVoxelsX; x++ {

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

func (grid *VoxelGrid) RectangleIntersects(rectCenter Vector3f, rectWidth, rectHeight int) bool {
	rectCorner := rectCenter.Sub(Vector3f{X: float32(rectWidth) / 2, Y: float32(rectHeight) / 2, Z: float32(rectWidth) / 2})
	for z := 0; z <= rectWidth; z++ {
		for y := 0; y <= rectHeight; y++ {
			for x := 0; x <= rectWidth; x++ {
				voxelPos := rectCorner.Plus(Vector3f{X: float32(x), Y: float32(y), Z: float32(z)}).DivScalar(grid.VoxelSize).ToVector3i()
				if grid.GetVoxel(voxelPos.X, voxelPos.Y, voxelPos.Z) {
					return true
				}
			}
		}
	}
	return false
}
