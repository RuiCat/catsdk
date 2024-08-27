// SPDX-License-Identifier: Unlicense OR MIT

/*
Package f32 is a float32 implementation of package image's
Point and affine transformations.

The coordinate space has the origin in the top left
corner with the axes extending right and down.
*/
package mat

import (
	"image"
	"math"
	"strconv"
)

// A Point is a two dimensional point.
type Point[T Float] struct {
	X, Y T
}

// String return a string representation of p.
func (p Point[T]) String() string {
	return "(" + strconv.FormatFloat(float64(p.X), 'f', -1, 32) +
		"," + strconv.FormatFloat(float64(p.Y), 'f', -1, 32) + ")"
}

// Pt is shorthand for Point[T]{X: x, Y: y}.
func Pt[T Float](x, y T) Point[T] {
	return Point[T]{X: x, Y: y}
}

// Add return the point p+p2.
func (p Point[T]) Add(p2 Point[T]) Point[T] {
	return Point[T]{X: p.X + p2.X, Y: p.Y + p2.Y}
}

// Sub returns the vector p-p2.
func (p Point[T]) Sub(p2 Point[T]) Point[T] {
	return Point[T]{X: p.X - p2.X, Y: p.Y - p2.Y}
}

// Mul returns p scaled by s.
func (p Point[T]) Mul(s T) Point[T] {
	return Point[T]{X: p.X * s, Y: p.Y * s}
}

// Div returns the vector p/s.
func (p Point[T]) Div(s T) Point[T] {
	return Point[T]{X: p.X / s, Y: p.Y / s}
}

// Round returns the integer point closest to p.
func (p Point[T]) Round() image.Point {
	return image.Point{
		X: int(math.Round(float64(p.X))),
		Y: int(math.Round(float64(p.Y))),
	}
}

// A Rectangle contains the points (X, Y) where Min.X <= X < Max.X,
// Min.Y <= Y < Max.Y.
type Rectangle[T Float] struct {
	Min, Max Point[T]
}

// String return a string representation of r.
func (r Rectangle[T]) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Rect is a shorthand for Rectangle[T]{Point[T]{x0, y0}, Point[T]{x1, y1}}.
// The returned Rectangle has x0 and y0 swapped if necessary so that
// it's correctly formed.
func Rect[T Float](x0, y0, x1, y1 T) Rectangle[T] {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle[T]{Point[T]{x0, y0}, Point[T]{x1, y1}}
}

// Size returns r's width and height.
func (r Rectangle[T]) Size() Point[T] {
	return Point[T]{X: r.Dx(), Y: r.Dy()}
}

// Dx returns r's width.
func (r Rectangle[T]) Dx() T {
	return r.Max.X - r.Min.X
}

// Dy returns r's Height.
func (r Rectangle[T]) Dy() T {
	return r.Max.Y - r.Min.Y
}

// Intersect returns the intersection of r and s.
func (r Rectangle[T]) Intersect(s Rectangle[T]) Rectangle[T] {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	if r.Empty() {
		return Rectangle[T]{}
	}
	return r
}

// Union returns the union of r and s.
func (r Rectangle[T]) Union(s Rectangle[T]) Rectangle[T] {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

// Canon returns the canonical version of r, where Min is to
// the upper left of Max.
func (r Rectangle[T]) Canon() Rectangle[T] {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// Empty reports whether r represents the empty area.
func (r Rectangle[T]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Add offsets r with the vector p.
func (r Rectangle[T]) Add(p Point[T]) Rectangle[T] {
	return Rectangle[T]{
		Point[T]{r.Min.X + p.X, r.Min.Y + p.Y},
		Point[T]{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub offsets r with the vector -p.
func (r Rectangle[T]) Sub(p Point[T]) Rectangle[T] {
	return Rectangle[T]{
		Point[T]{r.Min.X - p.X, r.Min.Y - p.Y},
		Point[T]{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Round returns the smallest integer rectangle that
// contains r.
func (r Rectangle[T]) Round() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(floor(r.Min.X)),
			Y: int(floor(r.Min.Y)),
		},
		Max: image.Point{
			X: int(ceil(r.Max.X)),
			Y: int(ceil(r.Max.Y)),
		},
	}
}

// fRect converts a rectangle to a f32internal.Rectangle.
func FRect[T Float](r image.Rectangle) Rectangle[T] {
	return Rectangle[T]{
		Min: FPt[T](r.Min), Max: FPt[T](r.Max),
	}
}

// Fpt converts an point to a f32.Point.
func FPt[T Float](p image.Point) Point[T] {
	return Point[T]{
		X: T(p.X), Y: T(p.Y),
	}
}

func ceil[T Float](v T) int {
	return int(math.Ceil(float64(v)))
}

func floor[T Float](v T) int {
	return int(math.Floor(float64(v)))
}
