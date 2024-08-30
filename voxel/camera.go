package voxel

import (
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
