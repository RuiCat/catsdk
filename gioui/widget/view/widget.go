package view

import (
	"gioui/op"
	"gioui/widget/material"
	"image"

	"gioui/widget/layout"
)

// OpsFace 上下文底层实现
type OpsFace interface {
	GetOps() *op.Ops           // 底层上下文
	GetTheme() *material.Theme // 底层主题
}

// LayoutFace 布局接口用于绘制组件
type LayoutFace interface {
	OpsFace                                       // 继承上下文
	Update()                                      // 状态更新
	Layout(gtx layout.Context) *layout.Dimensions // 组件布局
	GetPoint() *image.Point                       // 组件绘制位置
	GetDimensions() *layout.Dimensions            // 组件大小
	SetDisable(disable bool)                      // 组件禁用
	GetDisable() bool                             // 组件禁用
}

// UILayout 组件实现
type UILayout struct {
	LayoutFace
	Disable bool // 是否禁用组件
	*op.Ops
	*material.Theme    // 主题
	*image.Point       // 绘制位置
	*layout.Dimensions // 大小
}

func (ui *UILayout) GetOps() *op.Ops                   { return ui.Ops }
func (ui *UILayout) GetTheme() *material.Theme         { return ui.Theme }
func (ui *UILayout) GetPoint() *image.Point            { return ui.Point }
func (ui *UILayout) GetDimensions() *layout.Dimensions { return ui.Dimensions }
func (ui *UILayout) SetDisable(disable bool)           { ui.Disable = disable }
func (ui *UILayout) GetDisable() bool                  { return ui.Disable }
func (ui *UILayout) Update() {
	if ui.LayoutFace != nil {
		ui.LayoutFace.Update()
	}
}
func (ui *UILayout) Layout(gtx layout.Context) *layout.Dimensions { return ui.Dimensions }
