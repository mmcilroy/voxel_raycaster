package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

const WORLD_WIDTH, WORLD_HEIGHT = 256, 128

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

var raycaster = voxel.NewRaycastingCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66)

var sunPos voxel.Vector3f

var sunAngle float32

var sunHeight = float32(WORLD_HEIGHT - 1)

var enableLighting = true

var enableRecusiveDDA = true

var enableLOD = false

var enablePerPixelLighting = false

func initPerlinWorld(w, h int) *voxel.VoxelGrid {
	world := voxel.NewVoxelGrid(w, h, w, 1.0)

	perlinNoise := rl.GenImagePerlinNoise(world.NumVoxelsX, world.NumVoxelsZ, 0, 0, 0.5)
	colors := rl.LoadImageColors(perlinNoise)
	maxHeight := float32(0.0)
	gap := w / 6

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

func preUpdate() {
	dist := 1.3 * rl.GetFrameTime()

	if rl.IsKeyPressed('L') {
		enableLighting = !enableLighting
	}

	if rl.IsKeyPressed('R') {
		enableRecusiveDDA = !enableRecusiveDDA
	}

	if rl.IsKeyPressed('O') {
		enableLOD = !enableLOD
	}

	if rl.IsKeyPressed('P') {
		enablePerPixelLighting = !enablePerPixelLighting
	}

	if rl.IsKeyDown(rl.KeyLeftShift) {
		dist *= 20
	}

	if rl.IsKeyDown(rl.KeyDown) {
		sunHeight -= dist
	}

	if rl.IsKeyDown(rl.KeyUp) {
		sunHeight += dist
	}

	sunPos = scene.RotatingPosition(voxel.Vector3f{X: WORLD_WIDTH / 2, Y: WORLD_HEIGHT - 1, Z: WORLD_WIDTH / 2}, WORLD_WIDTH/2, sunAngle, 0)
	sunPos.Y = sunHeight
	sunAngle += rl.GetFrameTime()
}

func postUpdate() {
	rl.DrawText(fmt.Sprintf("Lighting (L): %t, RecursiveDDA (R): %t, LOD (O) %t, PerPixelLighting (P): %t", enableLighting, enableRecusiveDDA, enableLOD, enablePerPixelLighting), 20, 80, 20, rl.White)
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

func pixelColorFn(camera *voxel.RaycastingCamera, voxels *voxel.VoxelGrid, rayDir voxel.Vector3f) rl.Color {
	color := rl.SkyBlue

	// use the full resolution voxel grid if recursion is off
	// this should decrease performance
	if !enableRecusiveDDA {
		for voxels.Parent != nil {
			voxels = voxels.Parent
		}
	}

	var hit int
	var hitPos voxel.Vector3f
	var mapPos voxel.Vector3i

	if enableLOD {
		hit, hitPos, mapPos = voxels.DDARecursiveLOD(camera.Position, camera.Position, rayDir)
	} else {
		hit, hitPos, mapPos = voxels.DDARecursiveSimple(camera.Position, rayDir)
	}

	// if we are lighting per pixel do it per voxel face
	if !enablePerPixelLighting {
		hitPos = voxel.HitFaceCenter(hit, hitPos)
	}

	if hit != 0 {
		if enableLighting {
			// default unlit
			color = rl.Black

			// check if the hit point is visible to the sun
			sunHit, sunHitPos, sunMapPos := voxels.DDARecursiveSimple(sunPos, hitPos.Sub(sunPos).Normalize())

			// if sun ray hits the same block and face as our initial ray calc lighting
			if sunHit == hit && sunMapPos.Equals(mapPos) {
				diffuseLight := voxel.DiffuseLight(sunHit, voxel.Direction(sunPos, sunHitPos))
				color = pixelMinecraft(sunHit)
				color = rl.NewColor(
					uint8(float32(color.R)*diffuseLight),
					uint8(float32(color.G)*diffuseLight),
					uint8(float32(color.B)*diffuseLight),
					255)
			}
		} else {
			color = pixelMinecraft(hit)
		}
	}

	return color
}

func main() {
	// Full res world
	world := initPerlinWorld(WORLD_WIDTH, WORLD_HEIGHT)

	for world.NumVoxelsY > 2 {
		world = world.Compress()
	}

	scene.RenderRaycastingScene(&raycaster, world, pixelColorFn, preUpdate, postUpdate)
}
