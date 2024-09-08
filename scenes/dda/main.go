package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

const WORLD_WIDTH = 64
const WORLD_HEIGHT = 64
const VOXEL_SIZE = 1

var callbackCount = 0

var rayOrigin = voxel.Vector3fZero()

func readInput() {
	dist := 1.3 * rl.GetFrameTime()

	if rl.IsKeyDown(rl.KeyLeftShift) {
		dist *= 20
	}

	if rl.IsKeyDown('I') {
		rayOrigin.Z += dist
	}

	if rl.IsKeyDown('K') {
		rayOrigin.Z -= dist
	}

	if rl.IsKeyDown('J') {
		rayOrigin.X += dist
	}

	if rl.IsKeyDown('L') {
		rayOrigin.X -= dist
	}

	if rl.IsKeyDown('U') {
		rayOrigin.Y += dist
	}

	if rl.IsKeyDown('O') {
		rayOrigin.Y -= dist
	}
}

func ddaCallback(grid *voxel.VoxelGrid, x, y, z int) voxel.DDACallbackResult {
	scene.DrawVoxelOutline(x, y, z, grid.VoxelSize, rl.NewColor(255, 0, 0, 63))
	callbackCount += 1

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
}

func main() {
	// create world and put a single voxel in the middle
	vx, vy, vz := WORLD_WIDTH/2, WORLD_HEIGHT/2, WORLD_WIDTH/2
	voxels := voxel.NewVoxelGrid(WORLD_WIDTH, WORLD_HEIGHT, WORLD_WIDTH, VOXEL_SIZE)
	voxels.SetVoxel(vx, vy, vz, true)
	voxels = voxels.Compress() // 32
	voxels = voxels.Compress() // 16
	voxels = voxels.Compress() // 8
	voxels = voxels.Compress() // 4
	voxels = voxels.Compress() // 2

	rayEnd := voxel.Vector3f{X: WORLD_WIDTH/2 + 0.5, Y: WORLD_HEIGHT/2 + 0.5, Z: WORLD_WIDTH/2 + 0.5}

	render3D := func() {
		callbackCount = 0
		hit, _, _ := voxels.DDARecursive(rayOrigin, voxel.Direction(rayEnd, rayOrigin), ddaCallback)
		scene.DrawSphere(rayOrigin, 0.5, rl.Green)
		if hit != 0 {
			scene.DrawSphere(rayEnd, 0.5, rl.Red)
		}
		rl.DrawGrid(WORLD_WIDTH*2, 1)
	}

	render2D := func() {
		rl.DrawText(fmt.Sprintf("%d", callbackCount), 20, 20, 20, rl.Black)
	}

	scene.RenderScene(readInput, render3D, render2D)
}
