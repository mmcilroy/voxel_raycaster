package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

func main() {
	raycaster := voxel.NewRaycastingCamera(32, 18, 0.66)
	raycaster.Position = rl.Vector3{X: 1, Y: 0, Z: 1}
	rotX, rotY := float32(-5.5), float32(0.0)

	handleInput := func() {
		if rl.IsKeyDown('U') {
			rotX += 0.01
		}

		if rl.IsKeyDown('O') {
			rotX -= 0.01
		}

		if rl.IsKeyDown('J') {
			rotY += 0.01
		}

		if rl.IsKeyDown('L') {
			rotY -= 0.01
		}

		if rl.IsKeyDown('I') {
			raycaster.Direction.Z -= 0.01
		}

		if rl.IsKeyDown('K') {
			raycaster.Direction.Z += 0.01
		}
	}

	render3D := func() {
		raycaster.Rotate(rotX, rotY)

		p3, r3 := raycaster.GetRayForPixel(int32(raycaster.Resolution.X)-1, int32(raycaster.Resolution.Y)-1)
		p2, r2 := raycaster.GetRayForPixel(0, int32(raycaster.Resolution.Y)-1)
		p1, r1 := raycaster.GetRayForPixel(int32(raycaster.Resolution.X)-1, 0)
		p0, r0 := raycaster.GetRayForPixel(0, 0)

		rl.DrawSphere(p0, 0.1, rl.Red)
		rl.DrawSphere(p1, 0.1, rl.Red)
		rl.DrawSphere(p2, 0.1, rl.Red)
		rl.DrawSphere(p3, 0.1, rl.Red)
		rl.DrawSphere(rl.Vector3Add(raycaster.Position, raycaster.Forward), 0.1, rl.Red)

		rl.DrawRay(rl.NewRay(raycaster.Position, r0), rl.Black)
		rl.DrawRay(rl.NewRay(raycaster.Position, r1), rl.Black)
		rl.DrawRay(rl.NewRay(raycaster.Position, r2), rl.Black)
		rl.DrawRay(rl.NewRay(raycaster.Position, r3), rl.Black)
		rl.DrawRay(rl.NewRay(raycaster.Position, raycaster.Forward), rl.Black)

		rl.DrawSphere(raycaster.Position, 0.2, rl.Black)
		rl.DrawGrid(128, 1)
	}

	scene.RenderScene(handleInput, render3D, func() {})
}
