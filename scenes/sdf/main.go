package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/voxel_raycaster/scenes"
	"github.com/mmcilroy/voxel_raycaster/voxel"
)

var rayOrigin = voxel.Vector3f{X: 14.28, Y: 7.6, Z: 0.0}

var rayEnd = voxel.Vector3f{X: 4.53, Y: 14.96, Z: 10.64}

var displayDist float32

func readInput() {
	dist := 5 * rl.GetFrameTime()

	if rl.IsKeyDown(rl.KeyLeftShift) {
		dist *= 10
	}

	if rl.IsKeyDown('I') {
		rayOrigin.Z += dist
	}

	if rl.IsKeyDown('K') {
		rayOrigin.Z -= dist
	}

	if rl.IsKeyDown('J') {
		rayOrigin.X += dist
	}

	if rl.IsKeyDown('L') {
		rayOrigin.X -= dist
	}

	if rl.IsKeyDown('U') {
		rayOrigin.Y += dist
	}

	if rl.IsKeyDown('O') {
		rayOrigin.Y -= dist
	}
}

func initPerlinWorld(w, h int32) *voxel.VoxelGrid {
	world := voxel.NewVoxelGrid(w, h, w, 1.0)

	perlinNoise := rl.GenImagePerlinNoise(int(world.NumVoxelsX), int(world.NumVoxelsZ), 0, 0, 0.5)
	colors := rl.LoadImageColors(perlinNoise)
	maxHeight := float32(0.0)

	for z := int32(0); z < world.NumVoxelsZ; z++ {
		for x := int32(0); x < world.NumVoxelsX; x++ {
			color := colors[x+z*world.NumVoxelsX]
			height := float32(color.R) / 255.0 * float32(world.NumVoxelsY)
			if height > float32(maxHeight) {
				maxHeight = height
			}
			for y := int32(0); y < int32(height)+1; y++ {
				world.SetVoxel(x, y, z, true)
			}
		}
	}

	return world
}

func calcDistance(voxels *voxel.VoxelGrid, x, y, z int32) float32 {
	if voxels.GetVoxel(x, y, z) {
		return 0
	}

	v1 := voxel.Vector3f{X: float32(x), Y: float32(y), Z: float32(z)}
	minDist := float32(999999999.0)
	for iy := int32(0); iy < voxels.NumVoxelsY-1; iy++ {
		for iz := int32(0); iz < voxels.NumVoxelsZ-1; iz++ {
			for ix := int32(0); ix < voxels.NumVoxelsX-1; ix++ {
				if voxels.GetVoxel(ix, iy, iz) {
					v2 := voxel.Vector3f{X: float32(ix), Y: float32(iy), Z: float32(iz)}
					dist := voxel.Distance(v1, v2)
					if dist < minDist {
						minDist = dist
					}
				}
			}
		}
	}

	return minDist
}

func getDistance(voxels *voxel.VoxelGrid, sdf *[]float32, rayStart voxel.Vector3f) float32 {
	i := rayStart.ToVector3i()
	return (*sdf)[i.X+i.Y*voxels.NumVoxelsX+i.Z*voxels.NumVoxelsX*voxels.NumVoxelsY]
}

func trace(voxels *voxel.VoxelGrid, sdf *[]float32, rayStart, rayDir voxel.Vector3f) {
	for s := 0; s < 16; s++ {
		scene.DrawSphere(rayStart, 0.25, rl.Black)
		d := getDistance(voxels, sdf, rayStart)
		if d <= 0 {
			break
		}
		rl.DrawSphereWires(rl.NewVector3(rayStart.X, rayStart.Y, rayStart.Z), d, 10, 10, rl.Red)
		rayStart = rayStart.Plus(rayDir.MulScalar(d))
	}
}

func main() {
	voxels := initPerlinWorld(16, 16)

	sdf := make([]float32, voxels.NumVoxelsX*voxels.NumVoxelsY*voxels.NumVoxelsZ)
	minDist := float32(999999999.0)
	maxDist := float32(0.0)

	for y := int32(0); y < voxels.NumVoxelsY-1; y++ {
		fmt.Println(y)
		for z := int32(0); z < voxels.NumVoxelsZ-1; z++ {
			for x := int32(0); x < voxels.NumVoxelsX-1; x++ {
				dist := calcDistance(voxels, x, y, z)
				if dist < 0 {
					dist = 0
				}
				if dist < minDist {
					minDist = dist
				}
				if dist > maxDist {
					maxDist = dist
				}
				sdf[x+y*voxels.NumVoxelsX+z*voxels.NumVoxelsX*voxels.NumVoxelsY] = dist
			}
		}
	}

	scene.RenderScene(readInput, func() {

		for y := int32(0); y < voxels.NumVoxelsY-1; y++ {
			for z := int32(0); z < voxels.NumVoxelsZ-1; z++ {
				for x := int32(0); x < voxels.NumVoxelsX-1; x++ {
					if voxels.GetVoxel(x, y, z) {
						scene.DrawVoxel(x, y, z, 1, rl.NewColor(0, 255, 0, 255))
						scene.DrawVoxelOutline(x, y, z, 1, rl.Black)
					}
				}
			}
		}

		//trace(voxels, &sdf, rayOrigin, voxel.Direction(rayEnd, rayOrigin))

		d := getDistance(voxels, &sdf, rayOrigin)
		scene.DrawSphere(rayOrigin, 0.25, rl.Black)
		rl.DrawSphereWires(rl.NewVector3(rayOrigin.X, rayOrigin.Y, rayOrigin.Z), d, 10, 10, rl.Red)

	}, func() {
		i := rayOrigin.ToVector3i()
		rl.DrawText(fmt.Sprintf("%.02f, %.02f, %.02f", rayOrigin.X, rayOrigin.Y, rayOrigin.Z), 20, 40, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("%d, %d, %d", i.X, i.Y, i.Z), 20, 60, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("%.02f", displayDist), 20, 80, 20, rl.Black)
	})
}
