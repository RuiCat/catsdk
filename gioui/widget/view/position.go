package view

import (
	"gioui/op"
	"gioui/widget/layout"
)

// UIPosition 组件偏移
type UIPosition struct {
	LayoutFace
}

// NewUIPosition 组件偏移
func NewUIPosition(face LayoutFace) *UIPosition {
	return &UIPosition{LayoutFace: face}
}

// Layout 绘制
func (ui *UIPosition) Layout(gtx layout.Context) *layout.Dimensions {
	defer op.Offset(*ui.GetPoint()).Push(ui.GetOps()).Pop()
	return ui.LayoutFace.Layout(gtx)
}
