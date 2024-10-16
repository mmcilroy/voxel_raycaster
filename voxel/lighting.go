package voxel

func HitNormal(hit int32) Vector3f {
	if hit == -1 {
		return Vector3f{X: 1, Y: 0, Z: 0}
	} else if hit == 1 {
		return Vector3f{X: -1, Y: 0, Z: 0}
	} else if hit == -2 {
		return Vector3f{X: 0, Y: 1, Z: 0}
	} else if hit == 2 {
		return Vector3f{X: 0, Y: -1, Z: 0}
	} else if hit == -3 {
		return Vector3f{X: 0, Y: 0, Z: 1}
	} else if hit == 3 {
		return Vector3f{X: 0, Y: 0, Z: -1}
	}
	return Vector3fZero()
}

func HitFaceCenter(hit int32, hitPos Vector3f, mapPos Vector3i, voxelSize float32) Vector3f {
	voxelSizeHalved := voxelSize / 2
	if hit == -1 || hit == 1 {
		return Vector3f{X: hitPos.X, Y: float32(mapPos.Y)*voxelSize + voxelSizeHalved, Z: float32(mapPos.Z)*voxelSize + voxelSizeHalved}
	} else if hit == -2 || hit == 2 {
		return Vector3f{X: float32(mapPos.X)*voxelSize + voxelSizeHalved, Y: hitPos.Y, Z: float32(mapPos.Z)*voxelSize + voxelSizeHalved}
	} else if hit == -3 || hit == 3 {
		return Vector3f{X: float32(mapPos.X)*voxelSize + voxelSizeHalved, Y: float32(mapPos.Y)*voxelSize + voxelSizeHalved, Z: hitPos.Z}
	}
	return Vector3fZero()
}

func DiffuseLight(hit int32, dir Vector3f) float32 {
	diffuseLight := HitNormal(hit).DotProduct(dir)
	if diffuseLight < 0.5 {
		diffuseLight = 0.5
	}
	return diffuseLight
}
