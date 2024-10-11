package view

import (
	"gioui/op"
	"gioui/widget/layout"
	"mat/mat/spatial/f32"
)

// UIScale 组件缩放实现
type UIScale struct {
	LayoutFace
	value        float32 // 记录缩放倍数
	f32.Affine2D         // 缩放矩阵
}

// NewUIScale 组件缩放实现
func NewUIScale(face LayoutFace, value float32) *UIScale {
	ui := &UIScale{LayoutFace: face}
	ui.SetScale(value)
	return ui
}

// SetScale 设置缩放值
func (ui *UIScale) SetScale(value float32) {
	ui.value = value
	ui.Update()
}

// Update 更新
func (ui *UIScale) Update() {
	ui.Affine2D = ui.Affine2D.Scale(f32.Point{X: 1, Y: 1}, f32.Point{X: ui.value, Y: ui.value})
	ui.LayoutFace.Update()
}

// GetScale 获取缩放值
func (ui *UIScale) GetScale() float32 {
	return ui.value
}

// Layout 绘制
func (ui *UIScale) Layout(gtx layout.Context) layout.Dimensions {
	op.Affine(ui.Affine2D).Push(ui.GetOps())
	return ui.LayoutFace.Layout(gtx)
}
