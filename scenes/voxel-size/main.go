package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

const WORLD_SIZE = 4

const VOXEL_SIZE = 1

var raycastingScene scene.RaycastingScene

func preUpdate() {
	dist := 5 * rl.GetFrameTime()

	if rl.IsKeyDown(rl.KeyLeftShift) {
		dist *= 5
	}

	if rl.IsKeyPressed(rl.KeyEqual) {
		raycastingScene.UncompressedVoxels.VoxelSize *= 2
	}

	if rl.IsKeyPressed(rl.KeyMinus) {
		raycastingScene.UncompressedVoxels.VoxelSize /= 2
	}
}

func postUpdate() {
	rl.DrawText(fmt.Sprintf("Size: %.02f", raycastingScene.UncompressedVoxels.VoxelSize), 20, 80, 20, rl.White)
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

func initWorld() *voxel.VoxelGrid {
	var world = voxel.NewVoxelGrid(WORLD_SIZE, WORLD_SIZE, WORLD_SIZE, VOXEL_SIZE)

	world.SetVoxel(0, 0, 0, true)
	world.SetVoxel(1, 1, 1, true)

	return world
}

func main() {
	raycastingScene = scene.RaycastingScene{
		Voxels:                 initWorld(),
		Camera:                 voxel.NewCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66),
		SunPos:                 voxel.Vector3f{X: WORLD_SIZE - 1, Y: WORLD_SIZE - 1, Z: 0},
		EnableRecursiveDDA:     false,
		EnableLighting:         false,
		EnablePerPixelLighting: true,
	}

	raycastingScene.Camera.Body.Position = voxel.Vector3f{X: -1, Y: 1, Z: -1}

	scene.RenderRaycastingScene(&raycastingScene, pixelColorFn, preUpdate, postUpdate)
}
