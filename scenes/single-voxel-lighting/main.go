package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

var gridSize = 4
var voxelSize = float32(4)
var sunPos = rl.NewVector3(15, 15, 0)

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

func handleInput() {
	if rl.IsKeyDown('J') {
		sunPos.X += 3.3 * rl.GetFrameTime()
	}

	if rl.IsKeyDown('L') {
		sunPos.X -= 3.3 * rl.GetFrameTime()
	}

	if rl.IsKeyDown('I') {
		sunPos.Y += 3.3 * rl.GetFrameTime()
	}

	if rl.IsKeyDown('K') {
		sunPos.Y -= 3.3 * rl.GetFrameTime()
	}
}

func main() {
	voxels := voxel.NewVoxelGrid(gridSize, gridSize, gridSize, voxelSize)
	voxels.SetVoxel(1, 1, 1, true)
	lookAt := rl.NewVector3(6, 6, 6)
	sunAngle := float32(0)

	render3D := func() {
		positions := []rl.Vector3{
			rl.NewVector3(6, 10, 6), // top
			rl.NewVector3(6, 2, 6),  // bottom
			rl.NewVector3(10, 6, 6), // left
			rl.NewVector3(2, 6, 6),  // right
			rl.NewVector3(6, 6, 2),  // front
			rl.NewVector3(6, 6, 10), // back
		}

		sunPos = rotatingPosition(rl.NewVector3(6, 15, 6), 8, sunAngle, 0)
		sunAngle += rl.GetFrameTime()

		for _, pos := range positions {
			dir := rl.Vector3Normalize(rl.Vector3Subtract(lookAt, pos))
			hit, hitPos, _ := voxels.DDASimple(pos, dir)

			if hit != 0 {
				// draw a ray towards the sun
				sunDir := rl.Vector3Normalize(rl.Vector3Subtract(hitPos, sunPos))
				rl.DrawRay(rl.NewRay(sunPos, sunDir), rl.SkyBlue)

				// check if the hit point is visible to the sun
				sunHit, sunHitPos, _ := voxels.DDASimple(sunPos, sunDir)

				// hit point color is yellow if visible, black if not
				if sunHit != 0 && sunHit == hit {

					// draw the normal
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

					rl.DrawSphere(sunHitPos, 0.5, rl.NewColor(uint8(255*diffuseLight), uint8(255*diffuseLight), uint8(255*diffuseLight), 255))
				}
			}
		}
	}

	render2D := func() {
		rl.DrawText(fmt.Sprintf("%.02f %.02f %.02f", sunPos.X, sunPos.Y, sunPos.Z), 20, 40, 20, rl.Black)
	}

	scene.RenderVoxelScene(voxels, handleInput, render3D, render2D)
}
