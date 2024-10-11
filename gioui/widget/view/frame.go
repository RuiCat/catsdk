package view

import (
	"gioui/op"
	"gioui/op/clip"
	"gioui/op/paint"
	"gioui/widget/layout"
	"image"
)

// UIFrame 边框
type UIFrame struct {
	LayoutFace
	width float32
	op    clip.Op
}

// NewUIFrame 边框
func NewUIFrame(face LayoutFace, width float32) *UIFrame {
	ui := &UIFrame{LayoutFace: face}
	ui.SetWidth(width)
	return ui
}

// UIFrame 设置边框大小
func (ui *UIFrame) SetWidth(width float32) {
	ui.width = width
	ui.Update()
}

// Update 更新
func (ui *UIFrame) Update() {
	dim := ui.GetDimensions()
	ui.op = clip.Stroke{Path: clip.RRect{Rect: image.Rectangle{Max: dim.Size}}.Path(&op.Ops{}), Width: ui.width}.Op()
	ui.LayoutFace.Update()
}

// Layout 绘制
func (ui *UIFrame) Layout(gtx layout.Context) layout.Dimensions {
	theme := ui.GetTheme()
	paint.FillShape(ui.GetOps(), theme.Fg, ui.op)
	return ui.LayoutFace.Layout(gtx)
}
