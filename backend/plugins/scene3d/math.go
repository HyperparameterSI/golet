// Author: HyperparameterSI

package scene3d

import "math"

// Vec3 is a point or direction in a right-handed world; Y points up.
type Vec3 struct{ X, Y, Z float32 }

func (a Vec3) Add(b Vec3) Vec3      { return Vec3{a.X + b.X, a.Y + b.Y, a.Z + b.Z} }
func (a Vec3) Sub(b Vec3) Vec3      { return Vec3{a.X - b.X, a.Y - b.Y, a.Z - b.Z} }
func (a Vec3) Scale(s float32) Vec3 { return Vec3{a.X * s, a.Y * s, a.Z * s} }
func (a Vec3) Dot(b Vec3) float32   { return a.X*b.X + a.Y*b.Y + a.Z*b.Z }
func (a Vec3) Cross(b Vec3) Vec3 {
	return Vec3{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X}
}

// Unit returns zero for a zero-length vector.
func (a Vec3) Unit() Vec3 {
	length := float32(math.Sqrt(float64(a.Dot(a))))
	if length == 0 {
		return Vec3{}
	}
	return a.Scale(1 / length)
}
func finite(v float32) bool { return !math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0) }
func (a Vec3) valid() bool  { return finite(a.X) && finite(a.Y) && finite(a.Z) }

// Mat4 is a column-major matrix acting on column vectors. Use Identity for
// objects without a transform; Instance's zero matrix also means identity.
type Mat4 [16]float32

func Identity() Mat4          { return Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1} }
func Translation(v Vec3) Mat4 { m := Identity(); m[12], m[13], m[14] = v.X, v.Y, v.Z; return m }
func Scaling(v Vec3) Mat4     { m := Identity(); m[0], m[5], m[10] = v.X, v.Y, v.Z; return m }

// RotationY turns by radians about the up axis.
func RotationY(radians float32) Mat4 {
	s, c := float32(math.Sin(float64(radians))), float32(math.Cos(float64(radians)))
	m := Identity()
	m[0], m[2], m[8], m[10] = c, -s, s, c
	return m
}

// Mul returns a*b: b is applied first.
func (a Mat4) Mul(b Mat4) Mat4 {
	var out Mat4
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			for k := 0; k < 4; k++ {
				out[c*4+r] += a[k*4+r] * b[c*4+k]
			}
		}
	}
	return out
}
func (a Mat4) point(v Vec3) Vec3 {
	return Vec3{a[0]*v.X + a[4]*v.Y + a[8]*v.Z + a[12], a[1]*v.X + a[5]*v.Y + a[9]*v.Z + a[13], a[2]*v.X + a[6]*v.Y + a[10]*v.Z + a[14]}
}
func (a Mat4) affine() bool {
	for _, v := range a {
		if !finite(v) {
			return false
		}
	}
	return a[3] == 0 && a[7] == 0 && a[11] == 0 && a[15] == 1
}
