package voxel

type RaycastingCamera struct {
	Position    Vector3f // the camera location
	Direction   Vector3f // where camera is looking (z component is focal length)
	Forward     Vector3f
	Right       Vector3f
	Up          Vector3f // what way is up
	Resolution  Vector2f // image plane resolution (x is always 1)
	PlaneDelta  Vector2f // distanc+e to move along image plane for each ray increment
	AspectRatio float32
	RotateX     float32
	RotateY     float32
}

func NewRaycastingCamera(resX, resY int32, focalLength float32) RaycastingCamera {
	aspectRatio := float32(resX) / float32(resY)
	planeDeltaX := 1.0 / float32(resX-1)
	planeDeltaY := 1.0 / aspectRatio / float32(resY-1)
	return RaycastingCamera{
		Position:    Vector3fZero(),
		Direction:   Vector3f{X: 0, Y: 0, Z: float32(focalLength)},
		Up:          Vector3f{X: 0, Y: 1, Z: 0},
		Resolution:  Vector2f{X: float32(resX), Y: float32(resY)},
		PlaneDelta:  Vector2f{X: planeDeltaX, Y: planeDeltaY},
		AspectRatio: aspectRatio,
	}
}

func (camera *RaycastingCamera) Rotate(x, y float32) {
	camera.Forward = camera.Direction.RotateByAxisAngle(camera.Up, x).Normalize()
	camera.Right = camera.Forward.CrossProduct(camera.Up)
	camera.RotateX = x
	camera.RotateY = y
}

func (camera *RaycastingCamera) GetRayForPixel(x, y int32) (Vector3f, Vector3f) {
	planeX := camera.PlaneDelta.X * float32(x)
	planeY := camera.PlaneDelta.Y * float32(y)

	pos := camera.Direction.Plus(Vector3f{X: planeX - 0.5, Y: planeY - (1 / camera.AspectRatio / 2), Z: 0})
	pos = pos.RotateByAxisAngle(camera.Up, float32(camera.RotateX))
	pos = pos.RotateByAxisAngle(camera.Right, camera.RotateY)

	pos = pos.Plus(camera.Position)
	dir := pos.Sub(camera.Position).Normalize()

	return pos, dir
}
