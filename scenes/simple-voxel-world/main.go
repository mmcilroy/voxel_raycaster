package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/voxel_raycaster/scenes"
	"github.com/mmcilroy/voxel_raycaster/voxel"
)

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

const WORLD_SIZE = 16

const VOXEL_SIZE = 1

var raycastingScene scene.RaycastingScene

func preUpdate() {
	dist := 1.3 * rl.GetFrameTime()

	if rl.IsKeyDown(rl.KeyLeftShift) {
		dist *= 20
	}

	if rl.IsKeyPressed('1') {
		raycastingScene.SunPos = voxel.Vector3f{X: WORLD_SIZE - 1, Y: WORLD_SIZE - 1, Z: 0}
	}

	if rl.IsKeyPressed('2') {
		raycastingScene.SunPos = voxel.Vector3f{X: 0, Y: WORLD_SIZE - 1, Z: 0}
	}

	if rl.IsKeyPressed('3') {
		raycastingScene.SunPos = voxel.Vector3f{X: WORLD_SIZE - 1, Y: WORLD_SIZE - 1, Z: WORLD_SIZE - 1}
	}

	if rl.IsKeyPressed('4') {
		raycastingScene.SunPos = voxel.Vector3f{X: 0, Y: WORLD_SIZE - 1, Z: WORLD_SIZE - 1}
	}

	if rl.IsKeyPressed('5') {
		raycastingScene.SunPos = voxel.Vector3f{X: WORLD_SIZE / 2, Y: WORLD_SIZE - 1, Z: WORLD_SIZE / 2}
	}

	if rl.IsKeyDown(rl.KeyDown) {
		raycastingScene.SunPos.Y -= dist
	}

	if rl.IsKeyDown(rl.KeyUp) {
		raycastingScene.SunPos.Y += dist
	}
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

func column(world *voxel.VoxelGrid, x, y, z int32) {
	for h := int32(0); h < y; h++ {
		world.SetVoxel(x, h, z, true)
	}
}

func initWorld() *voxel.VoxelGrid {
	var world = voxel.NewVoxelGrid(WORLD_SIZE, WORLD_SIZE, WORLD_SIZE, VOXEL_SIZE)

	for z := int32(0); z < world.NumVoxelsZ; z++ {
		for x := int32(0); x < world.NumVoxelsX; x++ {
			world.SetVoxel(x, 0, z, true)
		}
	}

	center := int32(WORLD_SIZE / 2)

	column(world, center-1, 2, center+1)
	column(world, center+1, 2, center-1)
	column(world, center+1, 2, center+3)
	column(world, center+3, 2, center+1)

	column(world, center, 3, center+1)
	column(world, center+1, 3, center)
	column(world, center+1, 3, center+2)
	column(world, center+2, 3, center+1)

	column(world, center+1, 4, center+1)

	world.SetVoxel(0, 0, 0, false)

	return world
}

func main() {
	raycastingScene = scene.RaycastingScene{
		Voxels:                 initWorld(),
		Camera:                 voxel.NewCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66),
		SunPos:                 voxel.Vector3f{X: WORLD_SIZE - 1, Y: WORLD_SIZE - 1, Z: 0},
		EnableRecursiveDDA:     true,
		EnableLighting:         true,
		EnablePerPixelLighting: true,
	}

	raycastingScene.Camera.Body.Position = voxel.Vector3f{X: 1, Y: 2, Z: 1}
	raycastingScene.Camera.Body.Rotation.Y = -0.1

	scene.RenderRaycastingScene(&raycastingScene, pixelColorFn, preUpdate, func() {})
}
