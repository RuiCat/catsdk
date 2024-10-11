package view

import (
	"gioui/op"
	"gioui/widget"
	"gioui/widget/layout"
	"gioui/widget/material"
	"image"
)

// UIScroll 滚动条
type UIScroll struct {
	LayoutFace
	Use            bool                    // 滚动条生效
	Axis           layout.Axis             // 滚动条方向
	Distance       float32                 // 滚动条滚动距离
	Scrollbar      widget.Scrollbar        // 滚动条对象
	ScrollbarStyle material.ScrollbarStyle // 滚动条对象风格
	off            image.Point             // 绘制位置
	call           func(Distance float32)  // 位置回调
}

// NewUIScroll 创建滚动条
func NewUIScroll(face LayoutFace, Axis layout.Axis, call func(Distance float32)) *UIScroll {
	ui := &UIScroll{LayoutFace: face}
	ui.ScrollbarStyle = material.Scrollbar(ui.LayoutFace.GetTheme(), &ui.Scrollbar)
	ui.Use = true
	ui.call = call
	ui.SetAxis(Axis)
	return ui
}

// SetAxis 设置坐标轴
func (ui *UIScroll) SetAxis(Axis layout.Axis) {
	ui.Axis = Axis
	ui.Update()
}

// SetAxis 设置回调
func (ui *UIScroll) SetCall(call func(Distance float32)) {
	ui.call = call
}

// Update 更新
func (ui *UIScroll) Update() {
	dim := ui.GetDimensions()
	if ui.Axis == layout.Horizontal {
		ui.off = image.Point{Y: dim.Size.Y - 10}
	} else {
		ui.off = image.Point{X: dim.Size.X - 10}
	}
	ui.LayoutFace.Update()
}

// Layout 绘制
func (ui *UIScroll) Layout(gtx layout.Context) layout.Dimensions {
	if ui.Use && !ui.GetDisable() {
		dim := ui.GetDimensions()
		gtx.Ops = ui.GetOps()
		gtx.Constraints.Max.X = dim.Size.X - 5
		gtx.Constraints.Max.Y = dim.Size.Y - 5
		// 计算滚动条位置
		if ui.Scrollbar.Dragging() {
			ui.Distance += ui.Scrollbar.ScrollDistance()
			if ui.call != nil {
				ui.call(ui.Distance)
			}
		}
		TransformStack := op.Offset(ui.off).Push(gtx.Ops)
		ui.ScrollbarStyle.Layout(
			gtx,
			ui.Axis,
			ui.Distance,
			ui.Distance,
		)
		TransformStack.Pop()
	}
	return ui.LayoutFace.Layout(gtx)
}
