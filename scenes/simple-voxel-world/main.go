package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mmcilroy/structure_go/voxel"
)

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

const WORLD_SIZE = 16

const VOXEL_SIZE = 1

var raycaster = voxel.NewRaycastingCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66)

var world = voxel.NewVoxelGrid(WORLD_SIZE, WORLD_SIZE, WORLD_SIZE, VOXEL_SIZE)

var sunPos = rl.NewVector3(WORLD_SIZE-1, WORLD_SIZE-1, 0)

var rotationX, rotationY = float32(0.5), float32(-0.5)

var pixels = make([]rl.Color, NUM_RAYS_X*NUM_RAYS_Y)

func readInput() {
	mouseDelta := rl.GetMouseDelta()
	rotationX += mouseDelta.X * -0.003
	rotationY += mouseDelta.Y * -0.003

	dist := 1.3 * rl.GetFrameTime()

	if rl.IsKeyDown(rl.KeyLeftShift) {
		dist *= 5
	}

	if rl.IsKeyPressed('1') {
		sunPos = rl.NewVector3(WORLD_SIZE-1, WORLD_SIZE-1, 0)
	}

	if rl.IsKeyPressed('2') {
		sunPos = rl.NewVector3(0, WORLD_SIZE-1, 0)
	}

	if rl.IsKeyPressed('3') {
		sunPos = rl.NewVector3(WORLD_SIZE-1, WORLD_SIZE-1, WORLD_SIZE-1)
	}

	if rl.IsKeyPressed('4') {
		sunPos = rl.NewVector3(0, WORLD_SIZE-1, WORLD_SIZE-1)
	}

	if rl.IsKeyPressed('5') {
		sunPos = rl.NewVector3(WORLD_SIZE/2, WORLD_SIZE-1, WORLD_SIZE/2)
	}

	if rl.IsKeyDown(rl.KeyDown) {
		sunPos.Y -= dist
	}

	if rl.IsKeyDown(rl.KeyUp) {
		sunPos.Y += dist
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

	if rl.IsKeyDown('I') {
		sunPos.Z += dist
	}

	if rl.IsKeyDown('J') {
		sunPos.X += dist
	}

	if rl.IsKeyDown('K') {
		sunPos.Z -= dist
	}

	if rl.IsKeyDown('L') {
		sunPos.X -= dist
	}

	if rl.IsKeyDown('U') {
		sunPos.Y += dist
	}

	if rl.IsKeyDown('O') {
		sunPos.Y -= dist
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
	if rh == 1 {
		color = rl.DarkBrown
	} else if rh == 2 {
		color = rl.Green
	} else if rh == 3 {
		color = rl.Brown
	} else if rh == 4 {
		color = rl.Black
	}
	return color
}

func pixelMinecraftDiffuse(voxels *voxel.VoxelGrid, hit int, hitPos rl.Vector3, mapPos rl.Vector3) rl.Color {
	color := rl.SkyBlue

	if rl.Vector3Equals(mapPos, rl.NewVector3(float32(int(sunPos.X)), float32(int(sunPos.Y)), float32(int(sunPos.Z)))) {
		return rl.Yellow
	}

	if hit != 0 {
		// something was hit, so color will be at least black
		color = rl.Black

		// dont hit the sun
		world.SetVoxel(int(sunPos.X), int(sunPos.Y), int(sunPos.Z), false)

		// check if the hit point is visible to the sun
		sunDir := rl.Vector3Normalize(rl.Vector3Subtract(hitPos, sunPos))
		sunHit, sunHitPos, sunMapPos := voxels.DDASimple(sunPos, sunDir)
		world.SetVoxel(int(sunPos.X), int(sunPos.Y), int(sunPos.Z), true)

		// check the sun ray hit our block and on the same face as our initial ray
		if sunHit != 0 && sunHit == hit && rl.Vector3Equals(mapPos, sunMapPos) {

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

			color = rl.NewColor(uint8(255*diffuseLight), uint8(255*diffuseLight), uint8(255*diffuseLight), 255)
		}
	}

	return color
}

func raycast(world *voxel.VoxelGrid, xa, xb, ya, yb int32) {
	for y := ya; y < yb; y++ {
		for x := xa; x < xb; x++ {

			// Work out the ray direction
			_, rd := raycaster.GetRayForPixel(int32(x), int32(y))

			// Walk the ray until we hit a voxel
			rh, rp, mp := world.DDARecursive(raycaster.Position, rd, func(grid *voxel.VoxelGrid, x, y, z int) voxel.DDACallbackResult {
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
			pixels[x+y*NUM_RAYS_X] = pixelMinecraftDiffuse(world, rh, rp, mp)
		}
	}
}

func column(world *voxel.VoxelGrid, x, y, z int) {
	for h := 0; h < y; h++ {
		world.SetVoxel(x, h, z, true)
	}
}

func render(world *voxel.VoxelGrid, target rl.RenderTexture2D) {

	world.Clear()

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

	world.SetVoxel(int(sunPos.X), int(sunPos.Y), int(sunPos.Z), true)

	// Spread rays across threads
	raycast(world, 0, int32(NUM_RAYS_X), 0, int32(NUM_RAYS_Y))

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
	rl.DrawText(fmt.Sprintf("%.02f", world.VoxelSize), 20, 100, 20, rl.White)
}

func main() {
	raycaster.Position.Y = VOXEL_SIZE * 2

	rl.InitWindow(1600, 900, "")
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
