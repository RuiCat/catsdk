package draw

import "math"

// Axis 枚举类型，表示 3D 空间中的三个坐标轴
type Axis uint8

const (
	AxisNone Axis = iota // 表示没有指定坐标轴
	AxisX                // 表示 X 轴（水平方向）
	AxisY                // 表示 Y 轴（纵向方向）
	AxisZ                // 表示 Z 轴（深度方向）
)

// NewAxisDrawing 创建平面
// axis: 垂直轴
func NewAxisDrawing(axis Axis, box Box, color Color) *AxisDrawing {
	bx := box
	switch axis {
	case AxisX:
		bx.Min.X = 1
		bx.Max.X = 1
		bx.Min.Y = box.Min.Y
		bx.Min.Z = box.Min.X
		bx.Max.Y = box.Max.Y
		bx.Max.Z = box.Max.X
	case AxisY:
		bx.Min.Y = 1
		bx.Max.Y = 1
		bx.Min.X = box.Min.X
		bx.Min.Z = box.Min.Y
		bx.Max.X = box.Max.X
		bx.Max.Z = box.Max.Y
	case AxisZ:
		bx.Min.Z = 1
		bx.Max.Z = 1
		bx.Min.X = box.Min.X
		bx.Min.Y = box.Min.Y
		bx.Max.X = box.Max.X
		bx.Max.Y = box.Max.Y
	}
	return &AxisDrawing{Axis: axis, Box: bx, Material: Material{Color: color}}
}

// AxisDrawing 垂直坐标的绘图平面
type AxisDrawing struct {
	Box
	Axis     Axis // 坐标轴
	Material Material
}

func (a *AxisDrawing) Compile()                       {}
func (a *AxisDrawing) BoundingBox() Box               { return a.Box }
func (a *AxisDrawing) MaterialAt(vec Vector) Material { return a.Material }
func (a *AxisDrawing) UV(v Vector) Vector {
	v = v.Sub(a.Min).Div(a.Max.Sub(a.Min))
	switch a.Axis {
	case AxisX:
		return Vector{X: v.Z, Y: v.Y, Z: 0}
	case AxisY:
		return Vector{X: v.X, Y: v.Z, Z: 0}
	case AxisZ:
		return Vector{X: v.X, Y: v.Y, Z: 0}
	}
	return Vector{}
}
func (a *AxisDrawing) NormalAt(vec Vector) Vector {
	return Vector{X: 1, Y: 1, Z: 1}
}
func (a *AxisDrawing) Intersect(ray Ray) Hit {
	n := a.Min.Sub(ray.Origin).Div(ray.Direction)
	f := a.Max.Sub(ray.Origin).Div(ray.Direction)
	n, f = n.Min(f), n.Max(f)
	t0 := math.Max(math.Max(n.X, n.Y), n.Z)
	t1 := math.Min(math.Min(f.X, f.Y), f.Z)
	if t0 > 0 && t0 <= t1 {
		return Hit{a, t0, nil}
	}
	return NoHit
}
