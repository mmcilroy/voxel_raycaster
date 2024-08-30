package voxel

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func abs(x float32) float32 {
	return max(x, -x)
}

type DDACallbackResult int

const (
	HIT DDACallbackResult = iota
	MISS
	OOB
)

type DDACallbackFn func(grid *VoxelGrid, x, y, z int) DDACallbackResult

func (grid *VoxelGrid) DDA(worldPos rl.Vector3, rayDir rl.Vector3, callback DDACallbackFn) (int, rl.Vector3, rl.Vector3) {

	// convert worldPos to voxel space
	voxelPos := rl.Vector3Scale(worldPos, 1/grid.VoxelSize)

	// which box of the map we're in
	mapX, mapY, mapZ := int(voxelPos.X), int(voxelPos.Y), int(voxelPos.Z)

	// length of ray from one xyz side to next
	deltaDistX, deltaDistY, deltaDistZ := float32(1e30), float32(1e30), float32(1e30)

	if rayDir.X != 0 {
		deltaDistX = abs(1 / rayDir.X)
	}

	if rayDir.Y != 0 {
		deltaDistY = abs(1 / rayDir.Y)
	}

	if rayDir.Z != 0 {
		deltaDistZ = abs(1 / rayDir.Z)
	}

	// length of ray from current position to next x or y-side
	var sideDistX, sideDistY, sideDistZ float32

	// what direction to step in x or y-direction (either +1 or -1)
	var stepX, stepY, stepZ int

	// calculate step and initial sideDist
	if rayDir.X < 0 {
		stepX = -1
		sideDistX = (voxelPos.X - float32(mapX)) * deltaDistX
	} else {
		stepX = 1
		sideDistX = (float32(mapX) + 1.0 - voxelPos.X) * deltaDistX
	}

	if rayDir.Y < 0 {
		stepY = -1
		sideDistY = (voxelPos.Y - float32(mapY)) * deltaDistY
	} else {
		stepY = 1
		sideDistY = (float32(mapY) + 1.0 - voxelPos.Y) * deltaDistY
	}

	if rayDir.Z < 0 {
		stepZ = -1
		sideDistZ = (voxelPos.Z - float32(mapZ)) * deltaDistZ
	} else {
		stepZ = 1
		sideDistZ = (float32(mapZ) + 1.0 - voxelPos.Z) * deltaDistZ
	}

	hit, side := 0, 4
	dist := float32(0.0)

	for hit == 0 {

		result := callback(grid, mapX, mapY, mapZ)

		if result == OOB {
			hit = 0
			break
		}

		if result == HIT {
			hit = side
			break
		}

		//jump to next map square, either in x, y or z direction
		if sideDistX <= sideDistY && sideDistX <= sideDistZ {
			dist = sideDistX
			sideDistX += deltaDistX
			mapX += stepX
			side = 1 * stepX
		} else if sideDistY <= sideDistX && sideDistY <= sideDistZ {
			dist = sideDistY
			sideDistY += deltaDistY
			mapY += stepY
			side = 2 * stepY
		} else {
			dist = sideDistZ
			sideDistZ += deltaDistZ
			mapZ += stepZ
			side = 3 * stepZ
		}
	}

	return hit,
		rl.Vector3Add(worldPos, rl.Vector3Scale(rl.Vector3Scale(rayDir, dist), grid.VoxelSize)),
		rl.NewVector3(float32(mapX), float32(mapY), float32(mapZ))
}

func (grid *VoxelGrid) DDARecursive(rayPos rl.Vector3, rayDir rl.Vector3, callback DDACallbackFn) (int, rl.Vector3, rl.Vector3) {
	// For the max resolution grid we move rayStart back slightly
	// As this reduces the chances of it starting inside a voxel
	if grid.Parent == nil {
		rayPos = rl.Vector3Subtract(rayPos, rayDir)
	}

	// Perform the DDA
	hit, hitPos, mapPos := grid.DDA(rayPos, rayDir, callback)

	// Nothing was hit or there is no parent so return immediately
	if hit == 0 || grid.Parent == nil {
		return hit, hitPos, mapPos
	}

	// Something was hit and we have further checks we could do
	return grid.Parent.DDARecursive(hitPos, rayDir, callback)
}

func (grid *VoxelGrid) DDASimple(worldPos rl.Vector3, rayDir rl.Vector3) (int, rl.Vector3, rl.Vector3) {
	return grid.DDA(worldPos, rayDir, func(grid *VoxelGrid, x, y, z int) DDACallbackResult {
		if x < 0 || y < 0 || z < 0 {
			return OOB
		}

		if x >= grid.NumVoxelsX || y >= grid.NumVoxelsY || z >= grid.NumVoxelsZ {
			return OOB
		}

		if grid.GetVoxel(x, y, z) {
			return HIT
		}

		return MISS
	})
}
