package voxel

type RaycastResult int

const (
	RAYCAST_HIT RaycastResult = iota
	RAYCAST_MISS
	RAYCAST_OUT_OF_BOUNDS
	RAYCAST_MAX_STEPS
)

type RaycastCallback func(grid *VoxelGrid, mapPos Vector3i)

func (grid *VoxelGrid) isOutside(mapPos Vector3i) bool {
	// returns true if mapPos is outside the grid
	return mapPos.X < 0 ||
		mapPos.Y < 0 ||
		mapPos.Z < 0 ||
		mapPos.X >= grid.NumVoxelsX ||
		mapPos.Y >= grid.NumVoxelsY ||
		mapPos.Z >= grid.NumVoxelsZ
}

func (grid *VoxelGrid) isOutOfBounds(mapPos Vector3i, rayDir Vector3f) bool {
	// returns true if mapPos is outside grid and we are heading away from it
	return (mapPos.X < 0 && rayDir.X <= 0) ||
		(mapPos.Y < 0 && rayDir.Y <= 0) ||
		(mapPos.Z < 0 && rayDir.Z <= 0) ||
		(mapPos.X >= grid.NumVoxelsX && rayDir.X >= 0) ||
		(mapPos.Y >= grid.NumVoxelsY && rayDir.Y >= 0) ||
		(mapPos.Z >= grid.NumVoxelsZ && rayDir.Z >= 0)
}

func (grid *VoxelGrid) checkHit(mapPos Vector3i, rayDir Vector3f) RaycastResult {
	// check if the voxel is outside the grid
	if grid.isOutside(mapPos) {
		if grid.isOutOfBounds(mapPos, rayDir) {
			return RAYCAST_OUT_OF_BOUNDS
		} else {
			return RAYCAST_MISS
		}
	}

	// check if the current voxel is empty
	if grid.GetVoxel(mapPos.X, mapPos.Y, mapPos.Z) {
		return RAYCAST_HIT
	}

	return RAYCAST_MISS
}

func calcSideDist(rayPos Vector3f, rayDir Vector3f, deltaDist Vector3f, mapPos Vector3i) Vector3f {
	var sideDist Vector3f

	if rayDir.X < 0 {
		sideDist.X = (rayPos.X - float32(mapPos.X)) * deltaDist.X
	} else {
		sideDist.X = (float32(mapPos.X) + 1.0 - rayPos.X) * deltaDist.X
	}

	if rayDir.Y < 0 {
		sideDist.Y = (rayPos.Y - float32(mapPos.Y)) * deltaDist.Y
	} else {
		sideDist.Y = (float32(mapPos.Y) + 1.0 - rayPos.Y) * deltaDist.Y
	}

	if rayDir.Z < 0 {
		sideDist.Z = (rayPos.Z - float32(mapPos.Z)) * deltaDist.Z
	} else {
		sideDist.Z = (float32(mapPos.Z) + 1.0 - rayPos.Z) * deltaDist.Z
	}

	return sideDist
}

func (grid *VoxelGrid) Raycast(rayPos Vector3f, rayDir Vector3f) (int32, Vector3f, Vector3i) {
	return grid.RaycastC(rayPos, rayDir, nil)
}

func (grid *VoxelGrid) RaycastC(rayPos Vector3f, rayDir Vector3f, callback RaycastCallback) (int32, Vector3f, Vector3i) {
	// convert rayPos to voxel space
	rayPos = rayPos.DivScalar(grid.VoxelSize)

	// which box of the map we're in
	mapPos := rayPos.ToVector3i()

	// length of ray from one xyz side to next
	deltaDist := rayDir.Inverse().Abs()

	// what direction to step in x or y-direction (either +1 or -1)
	step := rayDir.Sign().ToVector3i()

	// length of ray from current position to next x or y-side
	sideDist := calcSideDist(rayPos, rayDir, deltaDist, mapPos)

	hit, side := int32(0), int32(4)
	dist := float32(0.0)

	for hit == 0 {
		// call optional callback
		if callback != nil {
			callback(grid, mapPos)
		}

		// check if we have hit anything
		result := grid.checkHit(mapPos, rayDir)

		// no point proceeding if OOB
		if result == RAYCAST_OUT_OF_BOUNDS {
			hit = 0
			break
		}

		// we hit something
		if result == RAYCAST_HIT {
			hit = side
			break
		}

		// jump to next map square, either in x, y or z direction
		if sideDist.X <= sideDist.Y && sideDist.X <= sideDist.Z {
			dist = sideDist.X
			sideDist.X += deltaDist.X
			mapPos.X += step.X
			side = 1 * step.X
		} else if sideDist.Y <= sideDist.X && sideDist.Y <= sideDist.Z {
			dist = sideDist.Y
			sideDist.Y += deltaDist.Y
			mapPos.Y += step.Y
			side = 2 * step.Y
		} else {
			dist = sideDist.Z
			sideDist.Z += deltaDist.Z
			mapPos.Z += step.Z
			side = 3 * step.Z
		}
	}

	hitPos := rayPos.Plus(rayDir.MulScalar(dist)).MulScalar(grid.VoxelSize)

	return hit, hitPos, mapPos
}

func (grid *VoxelGrid) RaycastRecursive(rayPos Vector3f, rayDir Vector3f) (int32, Vector3f, Vector3i) {
	return grid.RaycastRecursiveC(rayPos, rayDir, nil)
}

func (grid *VoxelGrid) RaycastRecursiveC(rayPos Vector3f, rayDir Vector3f, callback RaycastCallback) (int32, Vector3f, Vector3i) {
	// for the max resolution grid we move rayPos back slightly
	// as this reduces the chances of it starting inside a voxel
	if grid.Parent == nil {
		rayPos = rayPos.Sub(rayDir)
	}

	// perform the DDA
	hit, hitPos, mapPos := grid.RaycastC(rayPos, rayDir, callback)

	/*
		if hit == 5 {
			// nothing hit because we reached max steps
			// proceed using a lower res grid
			if grid.Child != nil {
				return grid.Child.RaycastRecursive(hitPos, rayDir, callback)
			} else {
				return grid.RaycastRecursive(hitPos, rayDir, callback)
			}
		}
	*/

	// nothing was hit or there is no parent so return immediately
	if hit == 0 || grid.Parent == nil {
		return hit, hitPos, mapPos
	}

	// something was hit
	// proceed using a high res grid
	return grid.Parent.RaycastRecursiveC(hitPos, rayDir, callback)
}
