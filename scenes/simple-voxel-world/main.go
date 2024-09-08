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

func preUpdate() {
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

func pixelMinecraftDiffuse(camera *voxel.RaycastingCamera, voxels *voxel.VoxelGrid, rayDir voxel.Vector3f) rl.Color {
	color := rl.SkyBlue

	hit, hitPos, mapPos := voxels.DDASimple(camera.Position, rayDir)

	hitPos = voxel.HitFaceCenter(hit, hitPos)

	if hit != 0 {
		// default unlit
		color = rl.Black

		// is the hit point visible to the sun
		sunDir := voxel.Direction(hitPos, sunPos)
		sunHit, sunHitPos, sunMapPos := voxels.DDASimple(sunPos, sunDir)

		// if visible calc diffuse light
		if sunHit == hit && sunMapPos.Equals(mapPos) {
			diffuseLight := voxel.DiffuseLight(sunHit, voxel.Direction(sunPos, sunHitPos))
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

	column(world, center-1, 2, center+1)
	column(world, center+1, 2, center-1)
	column(world, center+1, 2, center+3)
	column(world, center+3, 2, center+1)

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

	scene.RenderRaycastingScene(&raycaster, world, pixelMinecraftDiffuse, preUpdate, func() {})

}
