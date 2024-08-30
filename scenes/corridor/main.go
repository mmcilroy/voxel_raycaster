package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

const WORLD_SIZE = 16

const VOXEL_SIZE = 1

var raycaster = voxel.NewRaycastingCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66)

var world = voxel.NewVoxelGrid(WORLD_SIZE, WORLD_SIZE, WORLD_SIZE, VOXEL_SIZE)

func pixelColorFn(camera *voxel.RaycastingCamera, voxels *voxel.VoxelGrid, rayDir rl.Vector3) rl.Color {
	rayHit, _, _ := world.DDARecursive(raycaster.Position, rayDir, func(grid *voxel.VoxelGrid, x, y, z int) voxel.DDACallbackResult {
		if x < 0 || y < 0 || z < 0 {
			return voxel.OOB
		}

		if x >= grid.NumVoxelsX || y >= grid.NumVoxelsY || z >= grid.NumVoxelsZ {
			return voxel.OOB
		}

		if grid.GetVoxel(x, y, z) {
			return voxel.HIT
		}

		return voxel.MISS
	})

	color := rl.SkyBlue
	if rayHit == 1 || rayHit == -1 {
		color = rl.DarkBrown
	} else if rayHit == 2 || rayHit == -2 {
		color = rl.Green
	} else if rayHit == 3 || rayHit == -3 {
		color = rl.Brown
	} else if rayHit == 4 || rayHit == -4 {
		color = rl.Black
	}

	return color
}

func main() {
	raycaster.Position.Y = VOXEL_SIZE * 2

	for z := 0; z < world.NumVoxelsZ; z++ {
		for x := 0; x < world.NumVoxelsX; x++ {
			world.SetVoxel(x, 0, z, true)
		}
	}

	scene.RenderRaycastingScene(&raycaster, world, pixelColorFn)
}
