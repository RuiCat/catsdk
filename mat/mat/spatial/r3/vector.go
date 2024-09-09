// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"fmt"
	"mat/mat/spatial/s1"
	"math"
)

// Vec represents a point in ℝ³.
type Vec struct {
	X, Y, Z float64
}

type Vector = Vec

// ApproxEqual reports whether v and ov are equal within a small epsilon.
func (v Vec) ApproxEqual(ov Vec) bool {
	const epsilon = 1e-16
	return math.Abs(v.X-ov.X) < epsilon && math.Abs(v.Y-ov.Y) < epsilon && math.Abs(v.Z-ov.Z) < epsilon
}

func (v Vec) String() string { return fmt.Sprintf("(%0.24f, %0.24f, %0.24f)", v.X, v.Y, v.Z) }

// Norm returns the Vec's norm.
func (v Vec) Norm() float64 { return math.Sqrt(v.Dot(v)) }

// Norm2 returns the square of the norm.
func (v Vec) Norm2() float64 { return v.Dot(v) }

// Normalize returns a unit Vec in the same direction as v.
func (v Vec) Normalize() Vec {
	n2 := v.Norm2()
	if n2 == 0 {
		return Vec{0, 0, 0}
	}
	return v.MulScalar(1 / math.Sqrt(n2))
}

// IsUnit returns whether this Vec is of approximately unit length.
func (v Vec) IsUnit() bool {
	const epsilon = 5e-14
	return math.Abs(v.Norm2()-1) <= epsilon
}

// Abs returns the Vec with nonnegative components.
func (v Vec) Abs() Vec { return Vec{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)} }

// Add returns the standard Vec sum of v and ov.
func (v Vec) Add(ov Vec) Vec { return Vec{v.X + ov.X, v.Y + ov.Y, v.Z + ov.Z} }

// Sub returns the standard Vec difference of v and ov.
func (v Vec) Sub(ov Vec) Vec { return Vec{v.X - ov.X, v.Y - ov.Y, v.Z - ov.Z} }

func (a Vec) Mul(b Vec) Vec {
	return Vec{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}
func (a Vec) Div(b Vec) Vec {
	return Vec{a.X / b.X, a.Y / b.Y, a.Z / b.Z}
}

func (a Vec) Mod(b Vec) Vec {
	// as implemented in GLSL
	x := a.X - b.X*math.Floor(a.X/b.X)
	y := a.Y - b.Y*math.Floor(a.Y/b.Y)
	z := a.Z - b.Z*math.Floor(a.Z/b.Z)
	return Vec{x, y, z}
}

// Dot returns the standard dot product of v and ov.
func (v Vec) Dot(ov Vec) float64 {
	return float64(v.X*ov.X) + float64(v.Y*ov.Y) + float64(v.Z*ov.Z)
}

// Cross returns the standard cross product of v and ov.
func (v Vec) Cross(ov Vec) Vec {
	return Vec{
		float64(v.Y*ov.Z) - float64(v.Z*ov.Y),
		float64(v.Z*ov.X) - float64(v.X*ov.Z),
		float64(v.X*ov.Y) - float64(v.Y*ov.X),
	}
}

// Distance returns the Euclidean distance between v and ov.
func (v Vec) Distance(ov Vec) float64 { return v.Sub(ov).Norm() }

// Angle returns the angle between v and ov.
func (v Vec) Angle(ov Vec) s1.Angle {
	return s1.Angle(math.Atan2(v.Cross(ov).Norm(), v.Dot(ov))) * s1.Radian
}

func (a Vec) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

func (a Vec) LengthN(n float64) float64 {
	if n == 2 {
		return a.Length()
	}
	a = a.Abs()
	return math.Pow(math.Pow(a.X, n)+math.Pow(a.Y, n)+math.Pow(a.Z, n), 1/n)
}

func (a Vec) Negate() Vec {
	return Vec{-a.X, -a.Y, -a.Z}
}

func (a Vec) AddScalar(b float64) Vec {
	return Vec{a.X + b, a.Y + b, a.Z + b}
}

func (a Vec) SubScalar(b float64) Vec {
	return Vec{a.X - b, a.Y - b, a.Z - b}
}

func (a Vec) DivScalar(b float64) Vec {
	return Vec{a.X / b, a.Y / b, a.Z / b}
}

func (a Vec) MulScalar(b float64) Vec {
	return Vec{a.X * b, a.Y * b, a.Z * b}
}
func (a Vec) Min(b Vec) Vec {
	return Vec{math.Min(a.X, b.X), math.Min(a.Y, b.Y), math.Min(a.Z, b.Z)}
}

func (a Vec) Max(b Vec) Vec {
	return Vec{math.Max(a.X, b.X), math.Max(a.Y, b.Y), math.Max(a.Z, b.Z)}
}

func (a Vec) MinAxis() Vec {
	x, y, z := math.Abs(a.X), math.Abs(a.Y), math.Abs(a.Z)
	switch {
	case x <= y && x <= z:
		return Vec{1, 0, 0}
	case y <= x && y <= z:
		return Vec{0, 1, 0}
	}
	return Vec{0, 0, 1}
}

func (a Vec) MinComponent() float64 {
	return math.Min(math.Min(a.X, a.Y), a.Z)
}

func (a Vec) MaxComponent() float64 {
	return math.Max(math.Max(a.X, a.Y), a.Z)
}

func (n Vec) Reflect(i Vec) Vec {
	return i.Sub(n.MulScalar(2 * n.Dot(i)))
}

func (n Vec) Refract(i Vec, n1, n2 float64) Vec {
	nr := n1 / n2
	cosI := -n.Dot(i)
	sinT2 := nr * nr * (1 - cosI*cosI)
	if sinT2 > 1 {
		return Vec{}
	}
	cosT := math.Sqrt(1 - sinT2)
	return i.MulScalar(nr).Add(n.MulScalar(nr*cosI - cosT))
}

func (n Vec) Reflectance(i Vec, n1, n2 float64) float64 {
	nr := n1 / n2
	cosI := -n.Dot(i)
	sinT2 := nr * nr * (1 - cosI*cosI)
	if sinT2 > 1 {
		return 1
	}
	cosT := math.Sqrt(1 - sinT2)
	rOrth := (n1*cosI - n2*cosT) / (n1*cosI + n2*cosT)
	rPar := (n2*cosI - n1*cosT) / (n2*cosI + n1*cosT)
	return (rOrth*rOrth + rPar*rPar) / 2
}

// Axis enumerates the 3 axes of ℝ³.
type Axis int

// The three axes of ℝ³.
const (
	XAxis Axis = iota
	YAxis
	ZAxis
)

// Ortho returns a unit Vec that is orthogonal to v.
// Ortho(-v) = -Ortho(v) for all v.
func (v Vec) Ortho() Vec {
	ov := Vec{}
	switch v.LargestComponent() {
	case XAxis:
		ov.Z = 1
	case YAxis:
		ov.X = 1
	default:
		ov.Y = 1
	}
	return v.Cross(ov).Normalize()
}

// LargestComponent returns the axis that represents the largest component in this Vec.
func (v Vec) LargestComponent() Axis {
	t := v.Abs()

	if t.X > t.Y {
		if t.X > t.Z {
			return XAxis
		}
		return ZAxis
	}
	if t.Y > t.Z {
		return YAxis
	}
	return ZAxis
}

// SmallestComponent returns the axis that represents the smallest component in this Vec.
func (v Vec) SmallestComponent() Axis {
	t := v.Abs()

	if t.X < t.Y {
		if t.X < t.Z {
			return XAxis
		}
		return ZAxis
	}
	if t.Y < t.Z {
		return YAxis
	}
	return ZAxis
}

// Cmp compares v and ov lexicographically and returns:
//
//	-1 if v <  ov
//	 0 if v == ov
//	+1 if v >  ov
//
// This method is based on C++'s std::lexicographical_compare. Two entities
// are compared element by element with the given operator. The first mismatch
// defines which is less (or greater) than the other. If both have equivalent
// values they are lexicographically equal.
func (v Vec) Cmp(ov Vec) int {
	if v.X < ov.X {
		return -1
	}
	if v.X > ov.X {
		return 1
	}

	// First elements were the same, try the next.
	if v.Y < ov.Y {
		return -1
	}
	if v.Y > ov.Y {
		return 1
	}

	// Second elements were the same return the final compare.
	if v.Z < ov.Z {
		return -1
	}
	if v.Z > ov.Z {
		return 1
	}

	// Both are equal
	return 0
}

// Add returns the vector sum of p and q.
func Add(p, q Vec) Vec {
	return Vec{
		X: p.X + q.X,
		Y: p.Y + q.Y,
		Z: p.Z + q.Z,
	}
}

// Sub returns the vector sum of p and -q.
func Sub(p, q Vec) Vec {
	return Vec{
		X: p.X - q.X,
		Y: p.Y - q.Y,
		Z: p.Z - q.Z,
	}
}

// Scale returns the vector p scaled by f.
func Scale(f float64, p Vec) Vec {
	return Vec{
		X: f * p.X,
		Y: f * p.Y,
		Z: f * p.Z,
	}
}

// Dot returns the dot product p·q.
func Dot(p, q Vec) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

// Cross returns the cross product p×q.
func Cross(p, q Vec) Vec {
	return Vec{
		p.Y*q.Z - p.Z*q.Y,
		p.Z*q.X - p.X*q.Z,
		p.X*q.Y - p.Y*q.X,
	}
}

// Rotate returns a new vector, rotated by alpha around the provided axis.
func Rotate(p Vec, alpha float64, axis Vec) Vec {
	return NewRotation(alpha, axis).Rotate(p)
}

// Norm returns the Euclidean norm of p
//
//	|p| = sqrt(p_x^2 + p_y^2 + p_z^2).
func Norm(p Vec) float64 {
	return math.Hypot(p.X, math.Hypot(p.Y, p.Z))
}

// Norm2 returns the Euclidean squared norm of p
//
//	|p|^2 = p_x^2 + p_y^2 + p_z^2.
func Norm2(p Vec) float64 {
	return p.X*p.X + p.Y*p.Y + p.Z*p.Z
}

// Unit returns the unit vector colinear to p.
// Unit returns {NaN,NaN,NaN} for the zero vector.
func Unit(p Vec) Vec {
	if p.X == 0 && p.Y == 0 && p.Z == 0 {
		return Vec{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}
	}
	return Scale(1/Norm(p), p)
}

// Cos returns the cosine of the opening angle between p and q.
func Cos(p, q Vec) float64 {
	return Dot(p, q) / (Norm(p) * Norm(q))
}

// Divergence returns the divergence of the vector field at the point p,
// approximated using finite differences with the given step sizes.
func Divergence(p, step Vec, field func(Vec) Vec) float64 {
	sx := Vec{X: step.X}
	divx := (field(Add(p, sx)).X - field(Sub(p, sx)).X) / step.X
	sy := Vec{Y: step.Y}
	divy := (field(Add(p, sy)).Y - field(Sub(p, sy)).Y) / step.Y
	sz := Vec{Z: step.Z}
	divz := (field(Add(p, sz)).Z - field(Sub(p, sz)).Z) / step.Z
	return 0.5 * (divx + divy + divz)
}

// Gradient returns the gradient of the scalar field at the point p,
// approximated using finite differences with the given step sizes.
func Gradient(p, step Vec, field func(Vec) float64) Vec {
	dx := Vec{X: step.X}
	dy := Vec{Y: step.Y}
	dz := Vec{Z: step.Z}
	return Vec{
		X: (field(Add(p, dx)) - field(Sub(p, dx))) / (2 * step.X),
		Y: (field(Add(p, dy)) - field(Sub(p, dy))) / (2 * step.Y),
		Z: (field(Add(p, dz)) - field(Sub(p, dz))) / (2 * step.Z),
	}
}

// minElem return a vector with the minimum components of two vectors.
func minElem(a, b Vec) Vec {
	return Vec{
		X: math.Min(a.X, b.X),
		Y: math.Min(a.Y, b.Y),
		Z: math.Min(a.Z, b.Z),
	}
}

// maxElem return a vector with the maximum components of two vectors.
func maxElem(a, b Vec) Vec {
	return Vec{
		X: math.Max(a.X, b.X),
		Y: math.Max(a.Y, b.Y),
		Z: math.Max(a.Z, b.Z),
	}
}

// absElem returns the vector with components set to their absolute value.
func absElem(a Vec) Vec {
	return Vec{
		X: math.Abs(a.X),
		Y: math.Abs(a.Y),
		Z: math.Abs(a.Z),
	}
}

// mulElem returns the Hadamard product between vectors a and b.
//
//	v = {a.X*b.X, a.Y*b.Y, a.Z*b.Z}
func mulElem(a, b Vec) Vec {
	return Vec{
		X: a.X * b.X,
		Y: a.Y * b.Y,
		Z: a.Z * b.Z,
	}
}

// divElem returns the Hadamard product between vector a
// and the inverse components of vector b.
//
//	v = {a.X/b.X, a.Y/b.Y, a.Z/b.Z}
func divElem(a, b Vec) Vec {
	return Vec{
		X: a.X / b.X,
		Y: a.Y / b.Y,
		Z: a.Z / b.Z,
	}
}
