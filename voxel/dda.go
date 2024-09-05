package voxel

type DDACallbackResult int

const (
	HIT DDACallbackResult = iota
	MISS
	OOB
)

type DDACallbackFn func(grid *VoxelGrid, x, y, z int) DDACallbackResult

func (grid *VoxelGrid) DDA(rayPos Vector3f, rayDir Vector3f, callback DDACallbackFn) (int, Vector3f, Vector3i) {
	// convert rayPos to voxel space
	rayPos = rayPos.DivScalar(grid.VoxelSize)

	// which box of the map we're in
	mapPos := rayPos.ToVector3i()

	// length of ray from one xyz side to next
	deltaDist := rayDir.Inverse().Abs()

	// what direction to step in x or y-direction (either +1 or -1)
	step := rayDir.Sign().ToVector3i()

	// length of ray from current position to next x or y-side
	var sideDist Vector3f

	// calculate step and initial sideDist
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

	hit, side := 0, 4
	dist := float32(0.0)

	for hit == 0 {

		result := callback(grid, mapPos.X, mapPos.Y, mapPos.Z)

		if result == OOB {
			hit = 0
			break
		}

		if result == HIT {
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

func (grid *VoxelGrid) DDASimple(worldPos Vector3f, rayDir Vector3f) (int, Vector3f, Vector3i) {
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

func (grid *VoxelGrid) DDARecursive(rayPos Vector3f, rayDir Vector3f, callback DDACallbackFn) (int, Vector3f, Vector3i) {
	// For the max resolution grid we move rayStart back slightly
	// As this reduces the chances of it starting inside a voxel
	if grid.Parent == nil {
		rayPos = rayPos.Sub(rayDir)
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
