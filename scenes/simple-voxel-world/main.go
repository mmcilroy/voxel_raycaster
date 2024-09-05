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

var sunPos = voxel.Vector3f{X: WORLD_SIZE - 1, Y: WORLD_SIZE - 1, Z: 0}

func readInput() {
	dist := 1.3 * rl.GetFrameTime()

	if rl.IsKeyDown(rl.KeyLeftShift) {
		dist *= 20
	}

	if rl.IsKeyPressed('1') {
		sunPos = voxel.Vector3f{X: WORLD_SIZE - 1, Y: WORLD_SIZE - 1, Z: 0}
	}

	if rl.IsKeyPressed('2') {
		sunPos = voxel.Vector3f{X: 0, Y: WORLD_SIZE - 1, Z: 0}
	}

	if rl.IsKeyPressed('3') {
		sunPos = voxel.Vector3f{X: WORLD_SIZE - 1, Y: WORLD_SIZE - 1, Z: WORLD_SIZE - 1}
	}

	if rl.IsKeyPressed('4') {
		sunPos = voxel.Vector3f{X: 0, Y: WORLD_SIZE - 1, Z: WORLD_SIZE - 1}
	}

	if rl.IsKeyPressed('5') {
		sunPos = voxel.Vector3f{X: WORLD_SIZE / 2, Y: WORLD_SIZE - 1, Z: WORLD_SIZE / 2}
	}
	if rl.IsKeyDown(rl.KeyDown) {
		sunPos.Y -= dist
	}

	if rl.IsKeyDown(rl.KeyUp) {
		sunPos.Y += dist
	}
}

func pixelMinecraft(camera *voxel.RaycastingCamera, voxels *voxel.VoxelGrid, rayDir voxel.Vector3f) rl.Color {
	color := rl.SkyBlue
	hit, _, _ := voxels.DDASimple(camera.Position, rayDir)
	if hit == 1 || hit == -1 {
		color = rl.DarkBrown
	} else if hit == 2 || hit == -2 {
		color = rl.Green
	} else if hit == 3 || hit == -3 {
		color = rl.Brown
	} else if hit == 4 {
		color = rl.Black
	}
	return color
}

func pixelMinecraftDiffuse(camera *voxel.RaycastingCamera, voxels *voxel.VoxelGrid, rayDir voxel.Vector3f) rl.Color {
	color := rl.SkyBlue

	hit, hitPos, mapPos := voxels.DDASimple(camera.Position, rayDir)

	if hit != 0 {
		// something was hit, so color will be at least black
		color = rl.Black

		// check if the hit point is visible to the sun
		sunDir := hitPos.Sub(sunPos).Normalize()
		sunHit, sunHitPos, sunMapPos := voxels.DDASimple(sunPos, sunDir)

		// check the sun ray hit our block and on the same face as our initial ray
		if sunHit != 0 && sunHit == hit && mapPos.Equals(sunMapPos) {

			// calc normal
			normal := voxel.Vector3fZero()
			if sunHit == -1 {
				normal = voxel.Vector3f{X: 1, Y: 0, Z: 0}
			} else if sunHit == 1 {
				normal = voxel.Vector3f{X: -1, Y: 0, Z: 0}
			} else if sunHit == -2 {
				normal = voxel.Vector3f{X: 0, Y: 1, Z: 0}
			} else if sunHit == 2 {
				normal = voxel.Vector3f{X: 0, Y: -1, Z: 0}
			} else if sunHit == -3 {
				normal = voxel.Vector3f{X: 0, Y: 0, Z: 1}
			} else if sunHit == 3 {
				normal = voxel.Vector3f{X: 0, Y: 0, Z: -1}
			}

			lightDir := sunPos.Sub(sunHitPos).Normalize()
			diffuseLight := normal.DotProduct(lightDir)
			if diffuseLight < 0 {
				diffuseLight = 0
			}

			color = rl.NewColor(uint8(255*diffuseLight), uint8(255*diffuseLight), uint8(255*diffuseLight), 255)
		}
	}

	return color
}

func column(world *voxel.VoxelGrid, x, y, z int) {
	for h := 0; h < y; h++ {
		world.SetVoxel(x, h, z, true)
	}
}

func initWorld() *voxel.VoxelGrid {
	var world = voxel.NewVoxelGrid(WORLD_SIZE, WORLD_SIZE, WORLD_SIZE, VOXEL_SIZE)

	for z := 0; z < world.NumVoxelsZ; z++ {
		for x := 0; x < world.NumVoxelsX; x++ {
			world.SetVoxel(x, 0, z, true)
		}
	}

	center := WORLD_SIZE / 2

	// single blocks
	column(world, center, 2, center+2)
	column(world, center+2, 2, center)
	column(world, center+2, 2, center+3)
	column(world, center+3, 2, center+2)

	column(world, center, 3, center+1)
	column(world, center+1, 3, center)
	column(world, center+1, 3, center+2)
	column(world, center+2, 3, center+1)

	column(world, center+1, 4, center+1)

	return world
}

func main() {
	raycaster.Position.Y = VOXEL_SIZE * 2

	world := initWorld()

	scene.RenderRaycastingScene(&raycaster, world, pixelMinecraftDiffuse, readInput)

}
