package voxel

func HitNormal(hit int) Vector3f {
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

func HitFaceCenter(hit int, hitPos Vector3f) Vector3f {
	mapPos := hitPos.ToVector3i()
	if hit == -1 || hit == 1 {
		return Vector3f{X: hitPos.X, Y: float32(mapPos.Y) + 0.5, Z: float32(mapPos.Z) + 0.5}
	} else if hit == -2 || hit == 2 {
		return Vector3f{X: float32(mapPos.X) + 0.5, Y: hitPos.Y, Z: float32(mapPos.Z) + 0.5}
	} else if hit == -3 || hit == 3 {
		return Vector3f{X: float32(mapPos.X) + 0.5, Y: float32(mapPos.Y) + 0.5, Z: hitPos.Z}
	}
	return Vector3fZero()
}

func DiffuseLight(hit int, dir Vector3f) float32 {
	diffuseLight := HitNormal(hit).DotProduct(dir)
	if diffuseLight < 0 {
		diffuseLight = 0
	}
	return diffuseLight
}
