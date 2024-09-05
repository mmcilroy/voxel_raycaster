package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

var raycaster = voxel.NewRaycastingCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66)

var sunPos = voxel.Vector3f{X: 255, Y: 127, Z: 0}

var renderMode = 1

func initPerlinWorld(w, h int) *voxel.VoxelGrid {
	world := voxel.NewVoxelGrid(w, h, w, 1.0)

	perlinNoise := rl.GenImagePerlinNoise(world.NumVoxelsX, world.NumVoxelsZ, 0, 0, 0.5)
	colors := rl.LoadImageColors(perlinNoise)
	maxHeight := float32(0.0)
	gap := w / 8

	for z := gap; z < world.NumVoxelsZ-gap; z++ {
		for x := gap; x < world.NumVoxelsX-gap; x++ {
			color := colors[x+z*world.NumVoxelsX]
			height := float32(color.R) / 255.0 * float32(world.NumVoxelsY/2)
			if height > float32(maxHeight) {
				maxHeight = height
			}
			for y := 0; y < int(height)+1; y++ {
				world.SetVoxel(x, y, z, true)
			}
		}
	}

	raycaster.Position = voxel.Vector3f{X: 1, Y: maxHeight, Z: 1}

	return world
}

func readInput() {
	dist := 1.3 * rl.GetFrameTime()

	if rl.IsKeyDown(rl.KeyLeftShift) {
		dist *= 20
	}

	if rl.IsKeyPressed('1') {
		sunPos = voxel.Vector3f{X: 255, Y: 127, Z: 0}
	}

	if rl.IsKeyPressed('2') {
		sunPos = voxel.Vector3f{X: 0, Y: 127, Z: 0}
	}

	if rl.IsKeyPressed('3') {
		sunPos = voxel.Vector3f{X: 255, Y: 127, Z: 255}
	}

	if rl.IsKeyPressed('4') {
		sunPos = voxel.Vector3f{X: 0, Y: 127, Z: 255}
	}

	if rl.IsKeyPressed('5') {
		sunPos = voxel.Vector3f{X: 127, Y: 127, Z: 127}
	}

	if rl.IsKeyDown(rl.KeyDown) {
		sunPos.Y -= dist
	}

	if rl.IsKeyDown(rl.KeyUp) {
		sunPos.Y += dist
	}
}

func pixelMinecraft(rh int) rl.Color {
	color := rl.SkyBlue
	if rh == 1 || rh == -1 {
		color = rl.DarkBrown
	} else if rh == 2 || rh == -2 {
		color = rl.Green
	} else if rh == 3 || rh == -3 {
		color = rl.Brown
	} else if rh == 4 || rh == -4 {
		color = rl.Black
	}
	return color
}

func pixelColorFn(camera *voxel.RaycastingCamera, voxels *voxel.VoxelGrid, rayDir voxel.Vector3f) rl.Color {
	color := rl.SkyBlue

	hit, hitPos, mapPos := voxels.DDASimple(camera.Position, rayDir)

	if hit != 0 {
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

			color = pixelMinecraft(hit)
			color = rl.NewColor(
				uint8(float32(color.R)*diffuseLight),
				uint8(float32(color.G)*diffuseLight),
				uint8(float32(color.B)*diffuseLight),
				255)
		}
	}

	return color
}

func main() {
	// Full res world
	world := initPerlinWorld(256, 128)

	// Compress world here

	scene.RenderRaycastingScene(&raycaster, world, pixelColorFn, readInput)
}
