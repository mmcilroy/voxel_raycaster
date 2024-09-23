package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

var c voxel.Camera = voxel.NewCamera(1600, 900, 0.66)

func main() {
	handleInput := func() {
		var moveForward, moveSide, moveUp float32
		var rotateUp, rotateSide float32

		speed := 5 * rl.GetFrameTime()

		if rl.IsKeyDown('I') {
			moveForward += speed
		}

		if rl.IsKeyDown('K') {
			moveForward -= speed
		}

		if rl.IsKeyDown('J') {
			moveSide += speed
		}

		if rl.IsKeyDown('L') {
			moveSide -= speed
		}

		if rl.IsKeyDown('U') {
			if rl.IsKeyDown(rl.KeyLeftShift) {
				rotateUp += speed
			} else {
				rotateSide += speed
			}
		}

		if rl.IsKeyDown('O') {
			if rl.IsKeyDown(rl.KeyLeftShift) {
				rotateUp -= speed
			} else {
				rotateSide -= speed
			}
		}

		if rl.IsKeyDown(rl.KeySpace) {
			moveUp += speed
		}

		if rl.IsKeyDown(rl.KeyLeftControl) {
			moveUp -= speed
		}

		if rl.IsKeyDown('-') {
			c.FocalLength -= speed
		}

		if rl.IsKeyDown('=') {
			c.FocalLength += speed
		}

		c.Body.Move(moveForward, moveSide, moveUp)
		c.Body.Rotate(rotateSide, rotateUp)
	}

	render3D := func() {
		rl.DrawSphere(scene.ToRlVector(c.Body.Position), 0.2, rl.Black)
		rl.DrawRay(rl.NewRay(scene.ToRlVector(c.Body.Position), scene.ToRlVector(c.Body.Forward)), rl.Red)
		rl.DrawRay(rl.NewRay(scene.ToRlVector(c.Body.Position), scene.ToRlVector(c.Body.Right)), rl.Green)

		plane := c.Plane()
		rl.DrawSphere(scene.ToRlVector(plane.CenterPos), 0.2, rl.Black)
		rl.DrawRay(rl.NewRay(scene.ToRlVector(plane.CenterPos), scene.ToRlVector(plane.UpDir)), rl.Red)
		rl.DrawRay(rl.NewRay(scene.ToRlVector(plane.CenterPos), scene.ToRlVector(plane.RightDir)), rl.Green)

		rayPos, rayDir := c.RayDir(&plane, 0, 0)
		rl.DrawSphere(scene.ToRlVector(rayPos), 0.1, rl.Red)
		rl.DrawRay(rl.NewRay(scene.ToRlVector(rayPos), scene.ToRlVector(rayDir)), rl.Red)

		rayPos, rayDir = c.RayDir(&plane, c.Resolution.X-1, 0)
		rl.DrawSphere(scene.ToRlVector(rayPos), 0.1, rl.Green)
		rl.DrawRay(rl.NewRay(scene.ToRlVector(rayPos), scene.ToRlVector(rayDir)), rl.Green)

		rayPos, rayDir = c.RayDir(&plane, 0, c.Resolution.Y-1)
		rl.DrawSphere(scene.ToRlVector(rayPos), 0.1, rl.Blue)
		rl.DrawRay(rl.NewRay(scene.ToRlVector(rayPos), scene.ToRlVector(rayDir)), rl.Blue)

		rayPos, rayDir = c.RayDir(&plane, c.Resolution.X-1, c.Resolution.Y-1)
		rl.DrawSphere(scene.ToRlVector(rayPos), 0.1, rl.Yellow)
		rl.DrawRay(rl.NewRay(scene.ToRlVector(rayPos), scene.ToRlVector(rayDir)), rl.Yellow)

		rl.DrawGrid(128, 1)
	}

	scene.RenderScene(handleInput, render3D, func() {})
}
