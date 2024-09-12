package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mmcilroy/structure_go/voxel"
)

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

const WORLD_WIDTH = 512 // 64 meters

const WORLD_HEIGHT = 128 // 16 meters

const VOXEL_SIZE = 1.0 / 8.0

var world = voxel.NewVoxelGrid(WORLD_WIDTH, WORLD_HEIGHT, WORLD_WIDTH, VOXEL_SIZE)

var raycaster = voxel.NewRaycastingCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66)

var sunPos = voxel.Vector3f{X: 0, Y: (WORLD_HEIGHT - 1) * VOXEL_SIZE, Z: 0}

func preUpdate() {
	dist := 1.3 * rl.GetFrameTime()

	if rl.IsKeyDown(rl.KeyLeftShift) {
		dist *= 20
	}

	if rl.IsKeyDown(rl.KeyDown) {
		sunPos.Y -= dist
	}

	if rl.IsKeyDown(rl.KeyUp) {
		sunPos.Y += dist
	}
}

func pixelMinecraft(rh int) rl.Color {
	color := rl.Black
	if rh == 1 || rh == -1 {
		color = rl.Brown
	} else if rh == 2 || rh == -2 {
		color = rl.Green
	} else if rh == 3 || rh == -3 {
		color = rl.Brown
	}
	return color
}

func pixelColor(camera *voxel.RaycastingCamera, voxels *voxel.VoxelGrid, rayDir voxel.Vector3f) rl.Color {
	hit, _, _ := voxels.DDARecursiveSimple(camera.Position, rayDir)
	if hit != 0 {
		return pixelMinecraft(hit)
	} else {
		return rl.SkyBlue
	}
}

func initWorld() {
	/*
		for z := 0; z < world.NumVoxelsZ; z++ {
			for x := 0; x < world.NumVoxelsX; x++ {
				world.SetVoxel(x, 0, z, true)
			}
		}
	*/

	// 1
	world.SetVoxel(100, 0, 100, true)

	// 2x2
	world.SetVoxel(110, 0, 100, true)
	world.SetVoxel(111, 0, 100, true)
	world.SetVoxel(110, 1, 100, true)
	world.SetVoxel(111, 1, 100, true)

	// 4x4
	world.SetVoxel(120, 0, 100, true)
	world.SetVoxel(121, 0, 100, true)
	world.SetVoxel(122, 0, 100, true)
	world.SetVoxel(123, 0, 100, true)
	world.SetVoxel(120, 1, 100, true)
	world.SetVoxel(121, 1, 100, true)
	world.SetVoxel(122, 1, 100, true)
	world.SetVoxel(123, 1, 100, true)
	world.SetVoxel(120, 2, 100, true)
	world.SetVoxel(121, 2, 100, true)
	world.SetVoxel(122, 2, 100, true)
	world.SetVoxel(123, 2, 100, true)
	world.SetVoxel(120, 3, 100, true)
	world.SetVoxel(121, 3, 100, true)
	world.SetVoxel(122, 3, 100, true)
	world.SetVoxel(123, 3, 100, true)
}

func main() {
	raycaster.Position.X = 8
	raycaster.Position.Y = 8
	raycaster.Position.Z = 8

	initWorld()

	//scene.RenderRaycastingScene(&raycaster, world, pixelColor, preUpdate, func() {})
}
