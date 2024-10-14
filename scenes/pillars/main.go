package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/voxel_raycaster/scenes"
	"github.com/mmcilroy/voxel_raycaster/voxel"
)

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

const WORLD_SIZE = 64

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

func floor(world *voxel.VoxelGrid) {
	for z := int32(0); z < world.NumVoxelsZ; z++ {
		for x := int32(0); x < world.NumVoxelsX; x++ {
			world.SetVoxel(x, 0, z, true)
		}
	}
}

func column(world *voxel.VoxelGrid, x, y, z int32) {
	for h := int32(0); h < y; h++ {
		world.SetVoxel(x, h, z, true)
	}
}

func initWorld() *voxel.VoxelGrid {
	var world = voxel.NewVoxelGrid(WORLD_SIZE, WORLD_SIZE, WORLD_SIZE, VOXEL_SIZE)

	floor(world)

	column(world, WORLD_SIZE/2, WORLD_SIZE/2, WORLD_SIZE/2)
	column(world, 1+WORLD_SIZE/2, WORLD_SIZE/2, WORLD_SIZE/2)
	column(world, WORLD_SIZE/2, WORLD_SIZE/2, 1+WORLD_SIZE/2)
	column(world, 1+WORLD_SIZE/2, WORLD_SIZE/2, 1+WORLD_SIZE/2)

	return world
}

func main() {
	raycastingScene = scene.RaycastingScene{
		Voxels:                 initWorld(),
		Camera:                 voxel.NewCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66),
		SunPos:                 voxel.Vector3f{X: WORLD_SIZE, Y: WORLD_SIZE, Z: WORLD_SIZE},
		EnableRecursiveDDA:     true,
		EnableLighting:         true,
		EnablePerPixelLighting: true,
	}

	raycastingScene.Camera.Body.Position = voxel.Vector3f{X: 0, Y: 2, Z: 0}

	scene.RenderRaycastingScene(&raycastingScene, pixelColorFn, preUpdate, postUpdate)
}
