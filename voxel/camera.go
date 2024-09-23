package voxel

type Camera struct {
	Body        Moveable
	Resolution  Vector2i
	AspectRatio float32
	FocalLength float32
}

type CameraPlane struct {
	CenterPos Vector3f
	UpDir     Vector3f
	RightDir  Vector3f
}

func NewCamera(resX, resY, focalLength float32) Camera {
	return Camera{
		Resolution:  Vector2i{X: int32(resX), Y: int32(resY)},
		AspectRatio: resY / resX,
		FocalLength: focalLength,
	}
}

func (c *Camera) Plane() CameraPlane {
	plane := CameraPlane{}
	plane.RightDir = c.Body.Right
	plane.UpDir = UP.RotateByAxisAngle(c.Body.Right, c.Body.Rotation.Y).Normalize()
	plane.CenterPos = Vector3f{X: 0, Y: 0, Z: c.FocalLength}
	plane.CenterPos = plane.CenterPos.RotateByAxisAngle(UP, c.Body.Rotation.X)
	plane.CenterPos = plane.CenterPos.RotateByAxisAngle(c.Body.Right, c.Body.Rotation.Y)
	plane.CenterPos = plane.CenterPos.Plus(c.Body.Position)
	return plane
}

func (c *Camera) RayDir(plane *CameraPlane, x, y int32) (Vector3f, Vector3f) {
	rayPos := plane.CenterPos.Plus(plane.RightDir.MulScalar(float32(x)/float32(c.Resolution.X) - 0.5))
	rayPos = rayPos.Sub(plane.UpDir.MulScalar((float32(y)/float32(c.Resolution.Y) - 0.5) * c.AspectRatio))
	return rayPos, Direction(rayPos, c.Body.Position)
}
