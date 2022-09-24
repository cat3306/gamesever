package router

import (
	"fmt"
	"math"
)

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func (v *Vector3) String() string {
	return fmt.Sprintf("(%.2f, %.2f, %.2f)", v.X, v.Y, v.Z)
}

// DistanceTo calculates distance between two positions
func (v *Vector3) DistanceTo(o Vector3) float64 {
	dx := v.X - o.X
	dy := v.Y - o.Y
	dz := v.Z - o.Z
	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

// Sub calculates Vector3 p - Vector3 o
func (v *Vector3) Sub(o Vector3) Vector3 {
	return Vector3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

func (v *Vector3) Add(o Vector3) Vector3 {
	return Vector3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

// Mul calculates Vector3 p * m
func (v *Vector3) Mul(m float32) Vector3 {
	return Vector3{v.X * m, v.Y * m, v.Z * m}
}

// DirToYaw convert direction represented by Vector3 to Yaw
func (v *Vector3) DirToYaw() float32 {
	v.Normalize()

	yaw := math.Acos(float64(v.X))
	if v.Z < 0 {
		yaw = math.Pi*2 - yaw
	}

	yaw = yaw / math.Pi * 180 // convert to angle

	if yaw <= 90 {
		yaw = 90 - yaw
	} else {
		yaw = 90 + (360 - yaw)
	}

	return float32(yaw)
}

func (v *Vector3) Normalize() {
	d := float32(math.Sqrt(float64(v.X*v.X + v.Y + v.Y + v.Z*v.Z)))
	if d == 0 {
		return
	}
	v.X /= d
	v.Y /= d
	v.Z /= d
}

func (v *Vector3) Normalized() *Vector3 {
	v.Normalize()
	return v
}
