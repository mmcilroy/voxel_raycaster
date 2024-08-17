package main

import rl "github.com/gen2brain/raylib-go/raylib"

func abs(x float32) float32 {
	return max(x, -x)
}

func (grid *VoxelGrid) DDA(rayStart rl.Vector3, rayDir rl.Vector3) (int, rl.Vector3, rl.Vector3) {

	posX, posY, posZ := rayStart.X, rayStart.Y, rayStart.Z

	// which box of the map we're in
	mapX, mapY, mapZ := int(posX), int(posY), int(posZ)

	// length of ray from one x or y-side to next x or y-side
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

	//calculate step and initial sideDist
	if rayDir.X < 0 {
		stepX = -1
		sideDistX = (posX - float32(mapX)) * deltaDistX
	} else {
		stepX = 1
		sideDistX = (float32(mapX) + 1.0 - posX) * deltaDistX
	}

	if rayDir.Y < 0 {
		stepY = -1
		sideDistY = (posY - float32(mapY)) * deltaDistY
	} else {
		stepY = 1
		sideDistY = (float32(mapY) + 1.0 - posY) * deltaDistY
	}

	if rayDir.Z < 0 {
		stepZ = -1
		sideDistZ = (posZ - float32(mapZ)) * deltaDistZ
	} else {
		stepZ = 1
		sideDistZ = (float32(mapZ) + 1.0 - posZ) * deltaDistZ
	}

	hit, side := 0, 4
	dist := float32(0.0)

	for hit == 0 {
		// bounds check
		if mapX < 0 || mapY < 0 || mapZ < 0 {
			hit = 0
			break
		}

		if mapX >= grid.Size || mapY >= grid.Size || mapZ >= grid.Size {
			hit = 0
			break
		}

		// check current position for hit
		if grid.GetVoxel(mapX, mapY, mapZ) {
			hit = side
			break
		}

		//jump to next map square, either in x, y or z direction
		if sideDistX <= sideDistY && sideDistX <= sideDistZ {
			dist = sideDistX
			sideDistX += deltaDistX
			mapX += stepX
			side = 1
		} else if sideDistY <= sideDistX && sideDistY <= sideDistZ {
			dist = sideDistY
			sideDistY += deltaDistY
			mapY += stepY
			side = 2
		} else {
			dist = sideDistZ
			sideDistZ += deltaDistZ
			mapZ += stepZ
			side = 3
		}
	}

	return hit, rl.Vector3Add(rayStart, rl.Vector3Scale(rayDir, dist)), rl.NewVector3(float32(mapX), float32(mapY), float32(mapZ))
}

func (grid *VoxelGrid) DDA2Points(rayStart rl.Vector3, rayEnd rl.Vector3) (int, rl.Vector3, rl.Vector3) {
	rayDir := rl.Vector3Normalize(rl.Vector3Subtract(rayEnd, rayStart))
	hit, hitPos, mapPos := grid.DDA(rayStart, rayDir)

	// if there was nothing hit we return that
	if hit == 0 {
		return hit, hitPos, mapPos
	}

	// if something was hit then check it is before rayEnd
	pointDist := abs(rl.Vector3Distance(rayStart, rayEnd))
	hitDist := abs(rl.Vector3Distance(rayStart, hitPos))

	if hitDist < pointDist {
		return hit, hitPos, mapPos
	} else {
		return 0, rl.Vector3Zero(), rl.Vector3Zero()
	}
}

func (grid *VoxelGrid) DDARecursive(rayStart rl.Vector3, rayDir rl.Vector3, maxResolution int) (int, rl.Vector3, rl.Vector3) {
	// For the max resolution grid we move rayStart back slightly
	// As this reduces the chances of it starting inside a voxel
	if grid.Parent == nil {
		rayStart = rl.Vector3Subtract(rayStart, rayDir)
	}

	// Perform the DDA
	hit, hitPos, mapPos := grid.DDA(rayStart, rayDir)

	// Nothing was hit or there is no parent so return immediately
	if hit == 0 || grid.Parent == nil {
		return hit, hitPos, mapPos
	}

	// If the parent resolution is too high return
	if grid.Parent.Size > maxResolution {
		return hit, hitPos, mapPos
	}

	// Something was hit and we have further checks we could do
	return grid.Parent.DDARecursive(rl.Vector3Scale(hitPos, 4), rayDir, maxResolution)
}
