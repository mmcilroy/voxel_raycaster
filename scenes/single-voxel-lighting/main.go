package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

var sunPos = voxel.Vector3f{X: 15, Y: 15, Z: 0}

var height = float32(15.0)

func handleInput() {
	if rl.IsKeyDown('I') {
		height += 3.3 * rl.GetFrameTime()
	}

	if rl.IsKeyDown('K') {
		height -= 3.3 * rl.GetFrameTime()
	}
}

func main() {
	var gridSize = 4
	var voxelSize = float32(4)
	voxels := voxel.NewVoxelGrid(gridSize, gridSize, gridSize, voxelSize)
	voxels.SetVoxel(1, 1, 1, true)
	lookAt := voxel.Vector3f{X: 6, Y: 6, Z: 6}
	sunAngle := float32(0)

	render3D := func() {
		positions := []voxel.Vector3f{
			{X: 6, Y: 10, Z: 6}, // top
			{X: 6, Y: 2, Z: 6},  // bottom
			{X: 10, Y: 6, Z: 6}, // left
			{X: 2, Y: 6, Z: 6},  // right
			{X: 6, Y: 6, Z: 2},  // front
			{X: 6, Y: 6, Z: 10}, // back
		}

		sunPos = scene.RotatingPosition(voxel.Vector3f{X: 6, Y: 15, Z: 6}, 8, sunAngle, 0)
		sunPos.Y = height
		sunAngle += rl.GetFrameTime()

		for _, pos := range positions {
			// are we hitting the voxel
			dir := voxel.Direction(lookAt, pos)
			hit, hitPos, _ := voxels.DDASimple(pos, dir)

			// calc the light value for the hit
			if hit != 0 {
				// default unlit
				color := rl.Black

				// is the hit point visible to the sun
				sunDir := voxel.Direction(hitPos, sunPos)
				sunHit, sunHitPos, _ := voxels.DDASimple(sunPos, sunDir)

				// if visible calc diffuse light
				if sunHit == hit {
					diffuseLight := voxel.DiffuseLight(sunHit, voxel.Direction(sunPos, sunHitPos))
					color = rl.NewColor(uint8(255*diffuseLight), uint8(255*diffuseLight), uint8(255*diffuseLight), 255)
				}

				scene.DrawSphere(sunHitPos, 0.5, color)
				scene.DrawRay(sunPos, sunDir, rl.SkyBlue)
			}
		}
	}

	render2D := func() {
		rl.DrawText(fmt.Sprintf("%.02f %.02f %.02f", sunPos.X, sunPos.Y, sunPos.Z), 20, 40, 20, rl.Black)
	}

	scene.RenderVoxelScene(voxels, handleInput, render3D, render2D)
}
