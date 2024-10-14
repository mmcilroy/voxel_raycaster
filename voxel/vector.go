package voxel

import (
	"math"
)

type Vector2i struct {
	X, Y int32
}

type Vector3i struct {
	X, Y, Z int32
}

type Vector3f struct {
	X, Y, Z float32
}

type Vector2f struct {
	X, Y float32
}

func (v Vector3f) ToVector3i() Vector3i {
	return Vector3i{X: int32(v.X), Y: int32(v.Y), Z: int32(v.Z)}
}

func (v Vector3i) ToVector3f() Vector3f {
	return Vector3f{X: float32(v.X), Y: float32(v.Y), Z: float32(v.Z)}
}

func Vector2fZero() Vector2f {
	return Vector2f{X: 0, Y: 0}
}

func Vector3fZero() Vector3f {
	return Vector3f{X: 0, Y: 0, Z: 0}
}

func Vector3fUp() Vector3f {
	return Vector3f{X: 0, Y: 1, Z: 0}
}

func Vector3fRight() Vector3f {
	return Vector3f{X: -1, Y: 0, Z: 0}
}

func (v Vector3f) Sign() Vector3f {
	x, y, z := float32(1), float32(1), float32(1)
	if v.X < 0 {
		x = -1.0
	}
	if v.Y < 0 {
		y = -1.0
	}
	if v.Z < 0 {
		z = -1.0
	}
	return Vector3f{x, y, z}
}

func (v1 Vector3f) Plus(v2 Vector3f) Vector3f {
	return Vector3f{
		v1.X + v2.X,
		v1.Y + v2.Y,
		v1.Z + v2.Z,
	}
}

func (v1 Vector3f) PlusScalar(s float32) Vector3f {
	return Vector3f{
		v1.X + s,
		v1.Y + s,
		v1.Z + s,
	}
}

func (v1 Vector3f) Sub(v2 Vector3f) Vector3f {
	return Vector3f{
		v1.X - v2.X,
		v1.Y - v2.Y,
		v1.Z - v2.Z,
	}
}

func (v1 Vector3f) SubScalar(s float32) Vector3f {
	return Vector3f{
		v1.X - s,
		v1.Y - s,
		v1.Z - s,
	}
}

func (v1 Vector3f) Mul(v2 Vector3f) Vector3f {
	return Vector3f{
		v1.X * v2.X,
		v1.Y * v2.Y,
		v1.Z * v2.Z,
	}
}

func (v1 Vector3f) MulScalar(s float32) Vector3f {
	return Vector3f{
		v1.X * s,
		v1.Y * s,
		v1.Z * s,
	}
}

func (v1 Vector3f) Div(v2 Vector3f) Vector3f {
	return Vector3f{
		v1.X / v2.X,
		v1.Y / v2.Y,
		v1.Z / v2.Z,
	}
}

func (v1 Vector3f) DivScalar(s float32) Vector3f {
	return Vector3f{
		v1.X / s,
		v1.Y / s,
		v1.Z / s,
	}
}

func (v1 Vector3f) Inverse() Vector3f {
	return Vector3f{
		1 / v1.X,
		1 / v1.Y,
		1 / v1.Z,
	}
}

func (v1 Vector3f) Floor() Vector3f {
	return Vector3f{
		float32(math.Floor(float64(v1.X))),
		float32(math.Floor(float64(v1.Y))),
		float32(math.Floor(float64(v1.Z))),
	}
}

func (v1 Vector3f) Length() float32 {
	return float32(math.Sqrt(float64(v1.X*v1.X + v1.Y*v1.Y + v1.Z*v1.Z)))
}

func (v1 Vector3f) Abs() Vector3f {
	return Vector3f{
		max(v1.X, -v1.X),
		max(v1.Y, -v1.Y),
		max(v1.Z, -v1.Z),
	}
}

func (v1 Vector3f) RoundX() Vector3f {
	return Vector3f{
		float32(math.Round(float64(v1.X))),
		v1.Y,
		v1.Z,
	}
}

func (v1 Vector3f) RoundY() Vector3f {
	return Vector3f{
		v1.X,
		float32(math.Round(float64(v1.Y))),
		v1.Z,
	}
}

func (v1 Vector3f) RoundZ() Vector3f {
	return Vector3f{
		v1.X,
		v1.Y,
		float32(math.Round(float64(v1.Z))),
	}
}

func (v1 Vector3f) YZX() Vector3f {
	return Vector3f{
		v1.Y,
		v1.Z,
		v1.X,
	}
}

func (v1 Vector3f) ZXY() Vector3f {
	return Vector3f{
		v1.Z,
		v1.X,
		v1.Y,
	}
}

func (v1 Vector3f) Min(v2 Vector3f) Vector3f {
	return Vector3f{
		min(v1.X, v2.X),
		min(v1.Y, v2.Y),
		min(v1.Z, v2.Z),
	}
}

func (v1 Vector3f) LessThanEqual(v2 Vector3f) Vector3f {
	x, y, z := float32(0), float32(0), float32(0)
	if v1.X <= v2.X {
		x = 1.0
	}
	if v1.Y <= v2.Y {
		y = 1.0
	}
	if v1.Z <= v2.Z {
		z = 1.0
	}
	return Vector3f{X: x, Y: y, Z: z}
}

func (v Vector3f) Normalize() Vector3f {
	var length, ilength float32
	length = v.Length()
	if length == 0 {
		length = 1.0
	}
	ilength = 1.0 / length
	result := v
	result.X *= ilength
	result.Y *= ilength
	result.Z *= ilength
	return result
}

func (v Vector3f) RotateByAxisAngle(axis Vector3f, angle float32) Vector3f {
	// Using Euler-Rodrigues Formula
	// Ref.: https://en.wikipedia.org/w/index.php?title=Euler%E2%80%93Rodrigues_formula

	result := v

	// Vector3Normalize(axis);
	length := float32(math.Sqrt(float64(axis.X*axis.X + axis.Y*axis.Y + axis.Z*axis.Z)))
	if length == 0.0 {
		length = 1.0
	}
	ilength := 1.0 / length
	axis.X *= ilength
	axis.Y *= ilength
	axis.Z *= ilength

	angle /= 2.0
	a := float32(math.Sin(float64(angle)))
	b := axis.X * a
	c := axis.Y * a
	d := axis.Z * a
	a = float32(math.Cos(float64(angle)))
	w := Vector3f{X: b, Y: c, Z: d}

	// Vector3CrossProduct(w, v)
	wv := Vector3f{X: w.Y*v.Z - w.Z*v.Y, Y: w.Z*v.X - w.X*v.Z, Z: w.X*v.Y - w.Y*v.X}

	// Vector3CrossProduct(w, wv)
	wwv := Vector3f{X: w.Y*wv.Z - w.Z*wv.Y, Y: w.Z*wv.X - w.X*wv.Z, Z: w.X*wv.Y - w.Y*wv.X}

	// Vector3Scale(wv, 2*a)
	a *= 2
	wv.X *= a
	wv.Y *= a
	wv.Z *= a

	// Vector3Scale(wwv, 2)
	wwv.X *= 2
	wwv.Y *= 2
	wwv.Z *= 2

	result.X += wv.X
	result.Y += wv.Y
	result.Z += wv.Z

	result.X += wwv.X
	result.Y += wwv.Y
	result.Z += wwv.Z

	return result
}

func (v1 Vector3f) CrossProduct(v2 Vector3f) Vector3f {
	return Vector3f{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: v1.Z*v2.X - v1.X*v2.Z,
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

func (v1 Vector3f) DotProduct(v2 Vector3f) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vector3i) Equals(v2 Vector3i) bool {
	return v1.X == v2.X &&
		v1.Y == v2.Y &&
		v1.Z == v2.Z
}

func Direction(from Vector3f, to Vector3f) Vector3f {
	return from.Sub(to).Normalize()
}

func Distance(v1, v2 Vector3f) float32 {
	dx := v2.X - v1.X
	dy := v2.Y - v1.Y
	dz := v2.Z - v1.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}
