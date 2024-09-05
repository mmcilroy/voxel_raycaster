package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

func rotatingPosition(origin rl.Vector3, radius, angleX, angleY float32) rl.Vector3 {
	up := rl.NewVector3(0, 1, 0)
	pos := rl.NewVector3(0, 0, radius)
	pos = rl.Vector3RotateByAxisAngle(pos, up, angleX)
	forward := rl.Vector3Normalize(rl.Vector3RotateByAxisAngle(pos, up, angleX))
	right := rl.Vector3CrossProduct(forward, up)
	pos = rl.Vector3RotateByAxisAngle(pos, right, angleY)
	pos = rl.Vector3Add(pos, origin)
	return pos
}

func main() {
	mode := 1
	gridSize := 4
	voxelSize := float32(4)
	voxels := voxel.NewVoxelGrid(gridSize, gridSize, gridSize, voxelSize)
	voxels.SetVoxel(1, 1, 1, true)
	lookAt := rl.NewVector3(6, 6, 6)
	rotation := float32(0.0)

	readInput := func() {
		if rl.IsKeyPressed('1') {
			mode = 1
		}
		if rl.IsKeyPressed('2') {
			mode = 2
		}
		if rl.IsKeyPressed('3') {
			mode = 3
		}
	}

	render3D := func() {
		positions := []rl.Vector3{}
		if mode == 1 {
			positions = []rl.Vector3{
				rl.NewVector3(6, 10, 6), // top
				rl.NewVector3(6, 2, 6),  // bottom
				rl.NewVector3(10, 6, 6), // left
				rl.NewVector3(2, 6, 6),  // right
				rl.NewVector3(6, 6, 2),  // front
				rl.NewVector3(6, 6, 10), // back
			}
		} else if mode == 2 {
			positions = []rl.Vector3{
				rotatingPosition(rl.NewVector3(6, 6, 6), 6, rotation, 0), // rotating x
				rotatingPosition(rl.NewVector3(6, 6, 6), 6, 0, rotation), // rotating y
			}
		} else if mode == 3 {
			positions = []rl.Vector3{
				rotatingPosition(rl.NewVector3(6, 6, 6), 8, rotation, rotation), // rotating x and y
			}
		}

		for _, pos := range positions {
			dir := rl.Vector3Normalize(rl.Vector3Subtract(lookAt, pos))

			ro := voxel.Vector3f{X: pos.X, Y: pos.Y, Z: pos.Z}
			rd := voxel.Vector3f{X: dir.X, Y: dir.Y, Z: dir.Z}

			hit, hitPos, voxelPos := voxels.DDASimple(ro, rd)
			if hit != 0 {
				rl.DrawSphere(rl.NewVector3(hitPos.X, hitPos.Y, hitPos.Z), 0.2, rl.Yellow)
			} else {
				fmt.Println("Miss!")

			}
			if !rl.Vector3Equals(rl.NewVector3(float32(voxelPos.X), float32(voxelPos.Y), float32(voxelPos.Z)), rl.NewVector3(1, 1, 1)) {
				fmt.Println("VoxelPos!")
			}
			rl.DrawRay(rl.NewRay(pos, dir), rl.Black)
			rl.DrawSphere(pos, 0.2, rl.Green)
		}

		rotation += rl.GetFrameTime()
	}

	render2D := func() {
		rl.DrawText("Voxel Raycast Test. Press 1-3", 20, 20, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Mode: %d", mode), 20, 40, 20, rl.Black)
	}

	scene.RenderVoxelScene(voxels, readInput, render3D, render2D)
}
