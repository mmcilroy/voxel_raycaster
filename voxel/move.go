package voxel

var DEFAULT_DIRECTION = Vector3f{X: 0, Y: 0, Z: 1}
var UP = Vector3fUp()

type Moveable struct {
	Position Vector3f // current position
	Rotation Vector2f // how far we are looking to the left, right, up, down
	Forward  Vector3f // direction we will move forward or backward
	Right    Vector3f // direction we will move left or right
}

func (m *Moveable) Rotate(x, y float32) {
	m.Rotation.X += x
	m.Rotation.Y += y
	m.Forward = DEFAULT_DIRECTION.RotateByAxisAngle(UP, m.Rotation.X).Normalize()
	m.Right = m.Forward.CrossProduct(UP)
}

func (m *Moveable) Move(forward, right, up float32) {
	m.Position = m.Position.Plus(m.Forward.MulScalar(forward))
	m.Position = m.Position.Plus(m.Right.MulScalar(right))
	m.Position = m.Position.Plus(UP.MulScalar(up))
}
