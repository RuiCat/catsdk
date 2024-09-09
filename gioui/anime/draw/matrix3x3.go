package draw

import "math"

type Matrix3x3 struct {
	XX, YX, XY, YY, X0, Y0 float64
}

func Identity3x3() Matrix3x3 {
	return Matrix3x3{
		1, 0,
		0, 1,
		0, 0,
	}
}

func Translate3x3(x, y float64) Matrix3x3 {
	return Matrix3x3{
		1, 0,
		0, 1,
		x, y,
	}
}

func Scale3x3(x, y float64) Matrix3x3 {
	return Matrix3x3{
		x, 0,
		0, y,
		0, 0,
	}
}

func Rotate3x3(angle float64) Matrix3x3 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix3x3{
		c, s,
		-s, c,
		0, 0,
	}
}

func Shear3x3(x, y float64) Matrix3x3 {
	return Matrix3x3{
		1, y,
		x, 1,
		0, 0,
	}
}

func (a Matrix3x3) Multiply(b Matrix3x3) Matrix3x3 {
	return Matrix3x3{
		a.XX*b.XX + a.YX*b.XY,
		a.XX*b.YX + a.YX*b.YY,
		a.XY*b.XX + a.YY*b.XY,
		a.XY*b.YX + a.YY*b.YY,
		a.X0*b.XX + a.Y0*b.XY + b.X0,
		a.X0*b.YX + a.Y0*b.YY + b.Y0,
	}
}

func (a Matrix3x3) TransformVector(x, y float64) (tx, ty float64) {
	tx = a.XX*x + a.XY*y
	ty = a.YX*x + a.YY*y
	return
}

func (a Matrix3x3) TransformPoint(x, y float64) (tx, ty float64) {
	tx = a.XX*x + a.XY*y + a.X0
	ty = a.YX*x + a.YY*y + a.Y0
	return
}

func (a Matrix3x3) Translate(x, y float64) Matrix3x3 {
	return Translate3x3(x, y).Multiply(a)
}

func (a Matrix3x3) Scale(x, y float64) Matrix3x3 {
	return Scale3x3(x, y).Multiply(a)
}

func (a Matrix3x3) Rotate(angle float64) Matrix3x3 {
	return Rotate3x3(angle).Multiply(a)
}

func (a Matrix3x3) Shear(x, y float64) Matrix3x3 {
	return Shear3x3(x, y).Multiply(a)
}
