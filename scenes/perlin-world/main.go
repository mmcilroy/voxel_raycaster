package main

import (
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mmcilroy/structure_go/voxel"
)

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

var raycaster = voxel.NewRaycastingCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66)

var rotationX, rotationY = float32(-5.5), float32(-0.5)

var pixels = make([]rl.Color, NUM_RAYS_X*NUM_RAYS_Y)

var sunPos = rl.NewVector3(255, 127, 0)

var renderMode = 1

var wg sync.WaitGroup

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

	raycaster.Position = rl.Vector3{X: 1, Y: maxHeight, Z: 1}

	return world
}

func readInput() {
	mouseDelta := rl.GetMouseDelta()
	rotationX += mouseDelta.X * -0.003
	rotationY += mouseDelta.Y * -0.003

	dist := 1.3 * rl.GetFrameTime()

	if rl.IsKeyDown(rl.KeyLeftShift) {
		dist *= 20
	}

	if rl.IsKeyDown('W') {
		raycaster.Position = rl.Vector3Add(raycaster.Position, rl.Vector3Scale(raycaster.Forward, dist))
	}

	if rl.IsKeyDown('A') {
		raycaster.Position = rl.Vector3Subtract(raycaster.Position, rl.Vector3Scale(raycaster.Right, dist))
	}

	if rl.IsKeyDown('S') {
		raycaster.Position = rl.Vector3Subtract(raycaster.Position, rl.Vector3Scale(raycaster.Forward, dist))
	}

	if rl.IsKeyDown('D') {
		raycaster.Position = rl.Vector3Add(raycaster.Position, rl.Vector3Scale(raycaster.Right, dist))
	}

	if rl.IsKeyPressed('1') {
		sunPos = rl.NewVector3(255, 127, 0)
	}

	if rl.IsKeyPressed('2') {
		sunPos = rl.NewVector3(0, 127, 0)
	}

	if rl.IsKeyPressed('3') {
		sunPos = rl.NewVector3(255, 127, 255)
	}

	if rl.IsKeyPressed('4') {
		sunPos = rl.NewVector3(0, 127, 255)
	}

	if rl.IsKeyPressed('5') {
		sunPos = rl.NewVector3(127, 127, 127)
	}

	if rl.IsKeyDown(rl.KeyDown) {
		sunPos.Y -= dist
	}

	if rl.IsKeyDown(rl.KeyUp) {
		sunPos.Y += dist
	}

	if rl.IsKeyDown(rl.KeySpace) {
		raycaster.Position.Y += dist
	}

	if rl.IsKeyDown(rl.KeyLeftControl) {
		raycaster.Position.Y -= dist
	}

	raycaster.Rotate(rotationX, rotationY)
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

func pixelMinecraftDiffuse(voxels *voxel.VoxelGrid, hit int, hitPos rl.Vector3, mapPos rl.Vector3) rl.Color {
	color := rl.SkyBlue

	if hit != 0 {
		color = rl.Black

		// check if the hit point is visible to the sun
		sunDir := rl.Vector3Normalize(rl.Vector3Subtract(hitPos, sunPos))
		sunHit, sunHitPos, _ := voxels.DDASimple(sunPos, sunDir)

		// check the sun ray hit our block and on the same face as our initial ray
		if sunHit != 0 && sunHit == hit /*&& rl.Vector3Equals(mapPos, sunMapPos)*/ {

			// calc normal
			normal := rl.Vector3Zero()
			if sunHit == -1 {
				normal = rl.NewVector3(1, 0, 0)
			} else if sunHit == 1 {
				normal = rl.NewVector3(-1, 0, 0)
			} else if sunHit == -2 {
				normal = rl.NewVector3(0, 1, 0)
			} else if sunHit == 2 {
				normal = rl.NewVector3(0, -1, 0)
			} else if sunHit == -3 {
				normal = rl.NewVector3(0, 0, 1)
			} else if sunHit == 3 {
				normal = rl.NewVector3(0, 0, -1)
			}

			lightDir := rl.Vector3Normalize(rl.Vector3Subtract(sunPos, sunHitPos))
			diffuseLight := rl.Vector3DotProduct(normal, lightDir)
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

func raycast(world *voxel.VoxelGrid, xa, xb, ya, yb int32) {
	defer wg.Done()

	// Scale the ray start position to voxel space
	ro := raycaster.Position

	for y := ya; y < yb; y++ {
		for x := xa; x < xb; x++ {

			// Work out the ray direction
			_, rd := raycaster.GetRayForPixel(int32(x), int32(y))

			// Walk the ray until we hit a voxel
			rh, rp, mp := world.DDARecursive(ro, rd, func(grid *voxel.VoxelGrid, x, y, z int) voxel.DDACallbackResult {
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

			// Output the pixel color
			if renderMode == 0 {
				pixels[x+y*NUM_RAYS_X] = pixelMinecraft(rh)
			} else if renderMode == 1 {
				pixels[x+y*NUM_RAYS_X] = pixelMinecraftDiffuse(world, rh, rp, mp)
			}
		}
	}
}

const THREADS_NX = 4
const THREADS_NY = 4

func render(world *voxel.VoxelGrid, target rl.RenderTexture2D) {

	// Spread rays across threads
	wg.Add(THREADS_NX * THREADS_NY)
	for ty := 0; ty < THREADS_NY; ty++ {
		for tx := 0; tx < THREADS_NX; tx++ {
			go raycast(
				world,
				int32(NUM_RAYS_X/THREADS_NX*tx),
				int32(NUM_RAYS_X/THREADS_NX*(tx+1)),
				int32(NUM_RAYS_Y/THREADS_NY*ty),
				int32(NUM_RAYS_Y/THREADS_NY*(ty+1)))
		}
	}
	wg.Wait()

	// Use output color to create frame
	rl.BeginTextureMode(target)
	for ry := 0; ry < NUM_RAYS_Y; ry++ {
		for rx := 0; rx < NUM_RAYS_X; rx++ {
			rl.DrawPixel(NUM_RAYS_X-int32(rx), NUM_RAYS_Y-int32(ry), pixels[rx+ry*NUM_RAYS_X])
		}
	}
	rl.EndTextureMode()

	// Scale it
	rl.DrawTexturePro(target.Texture,
		rl.NewRectangle(0, 0, float32(target.Texture.Width), -float32(target.Texture.Height)),
		rl.NewRectangle(0, 0, 1600, 900),
		rl.NewVector2(0, 0),
		0,
		rl.White)

	rl.DrawFPS(20, 20)
	rl.DrawText(fmt.Sprintf("%.02f, %.02f, %.02f", raycaster.Position.X, raycaster.Position.Y, raycaster.Position.Z), 20, 40, 20, rl.White)
	rl.DrawText(fmt.Sprintf("%.02f, %.02f, %.02f", sunPos.X, sunPos.Y, sunPos.Z), 20, 60, 20, rl.White)
	rl.DrawText(fmt.Sprintf("%.02f, %.02f", rotationX, rotationY), 20, 80, 20, rl.White)
}

func main() {

	// Full res world
	world := initPerlinWorld(256, 128)

	// Compress world
	/*
		world = world.Compress()
		world = world.Compress()
		world = world.Compress()
	*/

	rl.InitWindow(1600, 900, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.DisableCursor()

	// Create a RenderTexture2D to use as a canvas
	target := rl.LoadRenderTexture(NUM_RAYS_X, NUM_RAYS_Y)

	// Clear render texture before entering the game loop
	rl.BeginTextureMode(target)
	rl.ClearBackground(rl.White)
	rl.EndTextureMode()

	for !rl.WindowShouldClose() {

		readInput()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		render(world, target)

		rl.EndDrawing()
	}
}
