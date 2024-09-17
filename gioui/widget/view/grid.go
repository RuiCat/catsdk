package view

import (
	"gioui/op"
	"gioui/op/clip"
	"gioui/op/paint"
	"gioui/widget/layout"
	"image"
	"mat/mat/spatial/f32"
)

// UIGrid 背景网格
type UIGrid struct {
	LayoutFace
	Use       bool        // 启用网格
	Size      image.Point // 网格大小
	OffsetPop *f32.Point  // 网格偏移
	distance  int         // 网格大小
	stroke    [4]clip.Op  // 组件背景网格与边框
	clipRect  clip.Op     // 用于组件剪裁
}

// NewUIGrid 背景网格
func NewUIGrid(face LayoutFace, distanceValue int, Size image.Point) *UIGrid {
	ui := &UIGrid{LayoutFace: face, Use: true, Size: Size}
	ui.SetDistance(distanceValue)
	return ui
}

// SetDistance 设置背景网格大小
func (ui *UIGrid) SetDistance(distanceValue int) {
	ui.distance = distanceValue
	ui.Update()
}

// GetDistance 获取缩放值
func (ui *UIGrid) GetDistance() int {
	return ui.distance
}

// Update 更新
func (ui *UIGrid) Update() {
	dim := ui.GetDimensions()
	x := dim.Size.X/ui.distance + 1
	y := dim.Size.Y/ui.distance + 1
	// 绘制背景网格
	var gridsX clip.Path
	gridsX.Begin(&op.Ops{})
	var gridsXZ clip.Path
	gridsXZ.Begin(&op.Ops{})
	for i := 0; i < x; i++ {
		ix := float32(i * ui.distance)
		gridsX.MoveTo(f32.Pt(ix, 0))
		gridsX.LineTo(f32.Pt(ix, float32(dim.Size.Y)))
		iz := float32((i + 1) * ui.distance)
		for ; ix < iz; ix += 10 {
			gridsXZ.MoveTo(f32.Pt(ix, 0))
			gridsXZ.LineTo(f32.Pt(ix, float32(dim.Size.Y)))
		}
	}
	var gridsY clip.Path
	gridsY.Begin(&op.Ops{})
	var gridsYZ clip.Path
	gridsYZ.Begin(&op.Ops{})
	for i := 0; i < y; i++ {
		iy := float32(i * ui.distance)
		gridsY.MoveTo(f32.Pt(0, iy))
		gridsY.LineTo(f32.Pt(float32(dim.Size.X), iy))
		iz := float32((i + 1) * ui.distance)
		for ; iy < iz; iy += 10 {
			gridsY.MoveTo(f32.Pt(0, iy))
			gridsY.LineTo(f32.Pt(float32(dim.Size.X), iy))
		}
	}
	ui.stroke[0] = clip.Stroke{Path: gridsX.End(), Width: 1.5}.Op()
	ui.stroke[1] = clip.Stroke{Path: gridsXZ.End(), Width: 0.5}.Op()
	ui.stroke[2] = clip.Stroke{Path: gridsY.End(), Width: 0.5}.Op()
	ui.stroke[3] = clip.Stroke{Path: gridsYZ.End(), Width: 1}.Op()
	ui.clipRect = clip.Rect{Max: ui.GetDimensions().Size}.Op()
	ui.LayoutFace.Update()
}

// Layout 绘制
func (ui *UIGrid) Layout(gtx layout.Context) *layout.Dimensions {
	defer ui.clipRect.Push(gtx.Ops).Pop()
	if ui.Use && !ui.GetDisable() {
		ops := ui.GetOps()
		theme := ui.GetTheme()
		// 绘制背景线
		if ui.OffsetPop != nil {
			off := image.Point{X: -int(ui.OffsetPop.X * float32(ui.Size.X)), Y: -int(ui.OffsetPop.Y * float32(ui.Size.Y))}
			offest := op.Offset(image.Point{X: off.X % ui.distance}).Push(gtx.Ops)
			paint.FillShape(ops, theme.Fg, ui.stroke[0])
			paint.FillShape(ops, theme.ContrastBg, ui.stroke[1])
			offest.Pop()
			offest = op.Offset(image.Point{Y: off.Y % ui.distance}).Push(gtx.Ops)
			paint.FillShape(ops, theme.Fg, ui.stroke[2])
			paint.FillShape(ops, theme.ContrastBg, ui.stroke[3])
			offest.Pop()
			defer op.Offset(off).Push(gtx.Ops).Pop()
		} else {
			paint.FillShape(ops, theme.Fg, ui.stroke[0])
			paint.FillShape(ops, theme.ContrastBg, ui.stroke[1])
			paint.FillShape(ops, theme.Fg, ui.stroke[2])
			paint.FillShape(ops, theme.ContrastBg, ui.stroke[3])
		}
	}
	return ui.LayoutFace.Layout(gtx)
}
