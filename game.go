package main

import (
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

var raycaster = NewRaycastingCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66)

var rotationX, rotationY = float32(-5.5), float32(-0.5)

var colors = make([]rl.Color, NUM_RAYS_X*NUM_RAYS_Y)

var voxelRes = 1024

var wg sync.WaitGroup

func initPerlinWorld(size int) *VoxelGrid {
	world := NewVoxelGrid(size)

	perlinNoise := rl.GenImagePerlinNoise(world.Size, world.Size, 0, 0, 0.5)
	colors := rl.LoadImageColors(perlinNoise)

	for z := 0; z < world.Size; z++ {
		for x := 0; x < world.Size; x++ {
			color := colors[x+z*world.Size]
			height := float32(color.R) / 255.0 * float32(world.Size) / 2.0
			for y := 0; y < int(height); y++ {
				world.SetVoxel(x, y, z, true)
			}
		}
	}

	raycaster.Position = rl.Vector3{X: 1, Y: 64, Z: 1}

	return world
}

func readInput() {
	mouseDelta := rl.GetMouseDelta()
	rotationX += mouseDelta.X * -0.003
	rotationY += mouseDelta.Y * -0.003

	dist := 1.3 * rl.GetFrameTime()

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

	if rl.IsKeyDown(rl.KeySpace) {
		raycaster.Position.Y += dist
	}

	if rl.IsKeyDown(rl.KeyLeftControl) {
		raycaster.Position.Y -= dist
	}

	if rl.IsKeyPressed(rl.KeyMinus) {
		voxelRes /= 4
	}

	if rl.IsKeyPressed(rl.KeyEqual) {
		voxelRes *= 4
	}

	raycaster.Rotate(rotationX, rotationY)
}

func raycast(world *VoxelGrid, xa, xb, ya, yb int32) {
	defer wg.Done()

	// Scale the ray start position to voxel space
	ro := raycaster.Position
	ro = rl.NewVector3(ro.X*world.Scale, ro.Y*world.Scale, ro.Z*world.Scale)

	for y := ya; y < yb; y++ {
		for x := xa; x < xb; x++ {

			// Work out the ray direction
			_, rd := raycaster.GetRayForPixel(int32(x), int32(y))

			// Walk the ray until we hit a voxel
			rh, _, _ := world.DDARecursive(ro, rd, voxelRes)

			// Choose the color for the pixel based on what was hit
			color := rl.SkyBlue
			if rh == 1 {
				color = rl.DarkGreen
			} else if rh == 2 {
				color = rl.Brown
			} else if rh == 3 {
				color = rl.Green
			} else if rh == 4 {
				color = rl.Black
			}

			// Output the pixel color
			colors[x+y*NUM_RAYS_X] = color
		}
	}
}

const THREADS_NX = 4
const THREADS_NY = 4

func render(world *VoxelGrid, target rl.RenderTexture2D) {

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
			rl.DrawPixel(NUM_RAYS_X-int32(rx), NUM_RAYS_Y-int32(ry), colors[rx+ry*NUM_RAYS_X])
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
	rl.DrawText(fmt.Sprintf("%.02f, %.02f, %.02f", raycaster.Position.X, raycaster.Position.Y, raycaster.Position.Z), 20, 40, 20, rl.Black)
	rl.DrawText(fmt.Sprintf("%.02f, %.02f", rotationX, rotationY), 20, 60, 20, rl.Black)
	rl.DrawText(fmt.Sprintf("Resolution: %d", voxelRes), 20, 80, 20, rl.Black)
}

func main() {

	world := initPerlinWorld(256)

	// Compress
	for world.Size > 4 {
		world = world.Compress()
		fmt.Printf("Compression level %d\n", world.Size)
	}

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
