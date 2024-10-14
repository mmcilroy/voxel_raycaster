package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/voxel_raycaster/scenes"
	"github.com/mmcilroy/voxel_raycaster/voxel"
)

const WORLD_WIDTH = 128
const WORLD_HEIGHT = 64
const VOXEL_SIZE = 1

var rayOrigin = voxel.Vector3fZero()
var rayEnd = voxel.Vector3f{X: WORLD_WIDTH/2 + 0.5, Y: WORLD_HEIGHT/2 + 0.5, Z: WORLD_WIDTH/2 + 0.5}

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

	if rl.IsKeyDown('T') {
		rayEnd.Z += dist
	}

	if rl.IsKeyDown('G') {
		rayEnd.Z -= dist
	}

	if rl.IsKeyDown('F') {
		rayEnd.X += dist
	}

	if rl.IsKeyDown('H') {
		rayEnd.X -= dist
	}

	if rl.IsKeyDown('R') {
		rayEnd.Y += dist
	}

	if rl.IsKeyDown('Y') {
		rayEnd.Y -= dist
	}
}

func main() {
	var numSteps int32

	// create world and put a single voxel in the middle
	vx, vy, vz := int32(WORLD_WIDTH/2), int32(WORLD_HEIGHT/2), int32(WORLD_WIDTH/2)
	voxels0 := voxel.NewTestVoxels(WORLD_WIDTH, WORLD_HEIGHT, WORLD_WIDTH, VOXEL_SIZE)
	voxels0.Set(vx, vy, vz, true)

	voxels1 := voxels0.Compress()    // 32
	voxels2 := (*voxels1).Compress() // 16
	voxels3 := (*voxels2).Compress() // 8
	voxels4 := (*voxels3).Compress() // 4

	tracer := voxel.MipmapTracerImpl{
		Voxels: []*voxel.Voxels{&voxels0, voxels1, voxels2, voxels3, voxels4},
	}

	render3D := func() {

		params := voxel.TraceParams{
			RayStart: rayOrigin,
			RayDir:   voxel.Direction(rayEnd, rayOrigin),
			MaxSteps: 0,
			Callback: func(voxels *voxel.Voxels, mapPos voxel.Vector3i) {
				scene.DrawVoxelOutline(mapPos.X, mapPos.Y, mapPos.Z, (*voxels).Size(), rl.NewColor(255, 0, 0, 63))
			},
		}

		result := tracer.Trace(params)
		numSteps = result.NumSteps

		scene.DrawVoxel(vx, vy, vz, 1, rl.Black)
		scene.DrawSphere(rayOrigin, 0.5, rl.Green)
		scene.DrawSphere(rayEnd, 0.5, rl.Blue)
		if result.Hit {
			scene.DrawSphere(result.HitPos, 0.5, rl.Red)
		}

		rl.DrawGrid(WORLD_WIDTH*2, 1)
	}

	render2D := func() {
		rl.DrawText(fmt.Sprintf("%d", numSteps), 20, 20, 20, rl.Black)
	}

	scene.RenderScene(readInput, render3D, render2D)
}
