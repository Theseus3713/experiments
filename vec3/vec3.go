package vec3

import (
	"math"
)

type Vector3 struct {
	X, Y, Z float32
}

func Add(v1, v2 Vector3) Vector3 {
	return Vector3{v1.X + v2.X,v1.Y + v2.Y,v1.Z + v2.Z}
}

func Mult(v Vector3, b float32) Vector3 {
	return Vector3{v.X * b, v.Y * b,v.Z * b}
}

func (v Vector3) Length() float32 {
	return float32(math.Sqrt(float64(v.X * v.X + v.Y * v.Y + v.Z * v.Z)))
}

func Distance(v1, v2 Vector3) float32 {
	var (
		xDiff = v1.X - v2.X
		yDiff = v1.Y - v2.Y
		zDiff = v1.Z - v2.Z
	)
	return float32(math.Sqrt(float64(xDiff * xDiff + yDiff * yDiff + zDiff * zDiff)))
}

func DistanceSquared(v1, v2 Vector3) float32 {
	var (
		xDiff = v1.X - v2.X
		yDiff = v1.Y - v2.Y
		zDiff = v1.Z - v2.Z
	)
	return xDiff * xDiff + yDiff * yDiff + zDiff * zDiff
}

func Normalize(v Vector3) Vector3 {
	len := v.Length()
	return Vector3{v.X / len, v.Y / len, v.Z / len}
}