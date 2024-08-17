package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaycastingCamera struct {
	Position    rl.Vector3 // the camera location
	Direction   rl.Vector3 // where camera is looking (z component is focal length)
	Forward     rl.Vector3
	Right       rl.Vector3
	Up          rl.Vector3 // what way is up
	Resolution  rl.Vector2 // image plane resolution (x is always 1)
	PlaneDelta  rl.Vector2 // distance to move along image plane for each ray increment
	AspectRatio float32
	RotateX     float32
	RotateY     float32
}

func NewRaycastingCamera(resX, resY int32, focalLength float32) RaycastingCamera {
	aspectRatio := float32(resX) / float32(resY)
	planeDeltaX := 1.0 / float32(resX-1)
	planeDeltaY := 1.0 / aspectRatio / float32(resY-1)
	return RaycastingCamera{
		Position:    rl.Vector3Zero(),
		Direction:   rl.NewVector3(0, 0, float32(focalLength)),
		Up:          rl.NewVector3(0, 1, 0),
		Resolution:  rl.NewVector2(float32(resX), float32(resY)),
		PlaneDelta:  rl.NewVector2(planeDeltaX, planeDeltaY),
		AspectRatio: aspectRatio,
	}
}

func (camera *RaycastingCamera) Rotate(x, y float32) {
	camera.Forward = rl.Vector3Normalize(rl.Vector3RotateByAxisAngle(camera.Direction, camera.Up, x))
	camera.Right = rl.Vector3CrossProduct(camera.Forward, camera.Up)
	camera.RotateX = x
	camera.RotateY = y
}

func (camera *RaycastingCamera) GetRayForPixel(x, y int32) (rl.Vector3, rl.Vector3) {
	planeX := camera.PlaneDelta.X * float32(x)
	planeY := camera.PlaneDelta.Y * float32(y)

	pos := rl.Vector3Add(camera.Direction, rl.NewVector3(planeX-0.5, planeY-(1/camera.AspectRatio/2), 0))
	pos = rl.Vector3RotateByAxisAngle(pos, camera.Up, float32(camera.RotateX))
	pos = rl.Vector3RotateByAxisAngle(pos, camera.Right, camera.RotateY)

	pos = rl.Vector3Add(pos, camera.Position)
	dir := rl.Vector3Normalize(rl.Vector3Subtract(pos, camera.Position))

	return pos, dir
}

func TestRaycastingCamera() {
	rl.InitWindow(1600, 900, "")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.DisableCursor()

	up := rl.Vector3{X: 0, Y: 1, Z: 0}
	camera := rl.Camera{
		Position:   rl.Vector3{X: 0, Y: 5, Z: -5},
		Target:     rl.Vector3{X: 0, Y: 0, Z: 0},
		Up:         up,
		Fovy:       60,
		Projection: rl.CameraPerspective,
	}

	raycaster := NewRaycastingCamera(32, 18, 0.66)
	raycaster.Position = rl.Vector3{X: 1, Y: 0, Z: 1}
	rotX, rotY := float32(-5.5), float32(0.0)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraThirdPerson)

		mouseDelta := rl.GetMouseDelta()

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

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		rl.BeginMode3D(camera)

		for z := 0; z < 64; z++ {
			for x := 0; x < 64; x++ {
				if x%2 == 0 && z%2 == 0 {
					rl.DrawCube(rl.NewVector3(float32(x), 0, float32(z)), 1, 1, 1, rl.NewColor(0, 255, 0, 64))
				}
			}
		}

		rl.DrawSphere(camera.Position, 0.5, rl.Red)

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
		rl.DrawSphere(rl.Vector3Zero(), 0.2, rl.Black)

		rl.DrawGrid(100, 1)

		rl.EndMode3D()

		rl.DrawFPS(10, 10)
		rl.DrawText(fmt.Sprintf("X: %.02f, Y: %.02f", mouseDelta.X, mouseDelta.Y), 10, 50, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Width: %.02f, Height: %.02f, Diagonal: %.02f,", rl.Vector3Distance(p0, p1), rl.Vector3Distance(p2, p3), rl.Vector3Distance(p3, p0)), 10, 100, 20, rl.Black)

		rl.EndDrawing()
	}
}
