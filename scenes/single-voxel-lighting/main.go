package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

var sunPos = voxel.Vector3f{X: 15, Y: 15, Z: 0}

var height = float32(15.0)

func rotatingPosition(origin voxel.Vector3f, radius, angleX, angleY float32) voxel.Vector3f {
	up := voxel.Vector3f{X: 0, Y: 1, Z: 0}
	pos := voxel.Vector3f{X: 0, Y: 0, Z: radius}
	pos = pos.RotateByAxisAngle(up, angleX)
	forward := pos.RotateByAxisAngle(up, angleX).Normalize()
	right := forward.CrossProduct(up)
	pos = pos.RotateByAxisAngle(right, angleY)
	pos = pos.Plus(origin)
	return pos
}

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

		sunPos = rotatingPosition(voxel.Vector3f{X: 6, Y: 15, Z: 6}, 8, sunAngle, 0)
		sunPos.Y = height
		sunAngle += rl.GetFrameTime()

		for _, pos := range positions {
			dir := voxel.Direction(lookAt, pos)
			hit, hitPos, _ := voxels.DDASimple(pos, dir)

			if hit != 0 {
				// draw a ray towards the sun
				sunDir := voxel.Direction(hitPos, sunPos)
				scene.DrawRay(sunPos, sunDir, rl.SkyBlue)

				// check if the hit point is visible to the sun
				sunHit, sunHitPos, _ := voxels.DDASimple(sunPos, sunDir)

				// hit point color is yellow if visible, black if not
				if sunHit != 0 && sunHit == hit {

					// draw the normal
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

					lightDir := voxel.Direction(sunPos, sunHitPos)
					diffuseLight := normal.DotProduct(lightDir)
					if diffuseLight < 0 {
						diffuseLight = 0
					}

					color := rl.NewColor(uint8(255*diffuseLight), uint8(255*diffuseLight), uint8(255*diffuseLight), 255)
					scene.DrawSphere(sunHitPos, 0.5, color)
				}
			}
		}
	}

	render2D := func() {
		rl.DrawText(fmt.Sprintf("%.02f %.02f %.02f", sunPos.X, sunPos.Y, sunPos.Z), 20, 40, 20, rl.Black)
	}

	scene.RenderVoxelScene(voxels, handleInput, render3D, render2D)
}
