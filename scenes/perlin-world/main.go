package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

const WORLD_WIDTH, WORLD_HEIGHT = 256, 128

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

func initPerlinWorld(w, h int32) *voxel.VoxelGrid {
	world := voxel.NewVoxelGrid(w, h, w, 1.0)

	perlinNoise := rl.GenImagePerlinNoise(int(world.NumVoxelsX), int(world.NumVoxelsZ), 0, 0, 0.5)
	colors := rl.LoadImageColors(perlinNoise)
	maxHeight := float32(0.0)
	gap := int32(0)

	for z := gap; z < world.NumVoxelsZ-gap; z++ {
		for x := gap; x < world.NumVoxelsX-gap; x++ {
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

func pixelColorFn(hit int32, mapPos voxel.Vector3i) rl.Color {
	color := rl.Black
	if hit == 1 || hit == -1 {
		color = rl.Brown
	} else if hit == 2 || hit == -2 {
		color = rl.Green
	} else if hit == 3 || hit == -3 {
		color = rl.Brown
	} else if hit == 0 {
		color = rl.SkyBlue
	}
	return color
}

func main() {
	world := initPerlinWorld(WORLD_WIDTH, WORLD_HEIGHT)

	raycastingScene := scene.RaycastingScene{
		Voxels:                 world,
		Camera:                 voxel.NewCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66),
		SunPos:                 voxel.Vector3f{X: WORLD_WIDTH - 1, Y: WORLD_HEIGHT - 1, Z: 0},
		EnableRecursiveDDA:     true,
		EnableLighting:         true,
		EnablePerPixelLighting: true,
	}

	raycastingScene.Camera.Body.Position = voxel.Vector3f{X: 16, Y: 96, Z: 16}

	scene.RenderRaycastingScene(&raycastingScene, pixelColorFn, func() {}, func() {})
}
