package view

import (
	"gioui/io/event"
	"gioui/io/pointer"
	"gioui/widget/layout"
	"mat/mat/spatial/f32"
)

// UIDrag 拖放
type UIDrag struct {
	LayoutFace
	Use          bool // 启用拖放
	dragging     bool
	position     f32.Point
	IsAdjustThe  UIDragAdjustThe // 大小调整
	IsAdjustMask UIDragAdjustThe
}

// NewUIDrag 创建拖动效果
func NewUIDrag(face LayoutFace) *UIDrag {
	ui := &UIDrag{LayoutFace: face, IsAdjustMask: 255}
	return ui
}

type UIDragAdjustThe uint8

const (
	UIDragTheW        UIDragAdjustThe = 1 << iota // 右
	UIDragTheE                                    // 左
	UIDragTheN                                    // 顶
	UIDragTheR                                    // 底
	UIDragTheRE                                   // 左下角
	UIDragTheDragging                             // 移动

)

// Layout 绘制
func (ui *UIDrag) Layout(gtx layout.Context) (dim layout.Dimensions) {
	event.Op(gtx.Ops, ui)
	if ev, ok := gtx.Source.Event(pointer.Filter{
		Target: ui,
		Kinds:  pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Cancel,
	}); ok {
		dim := ui.GetDimensions()
		point := ui.GetPoint()
		event := ev.(pointer.Event)
		pos := event.Position.Sub(f32.Point{X: float32(point.X), Y: float32(point.Y)})
		if event.Buttons == pointer.ButtonSecondary {
			ui.Use = true
			ui.SetDisable(true)
		} else if ((pos.X > 0 && pos.Y > 0 && pos.X < float32(dim.Size.X) && pos.Y < float32(dim.Size.Y)) && ui.Use) || ui.dragging {
			if !ui.dragging {
				//@ 处理鼠标所在点击的位置
				if pos.X+10 > float32(dim.Size.X) && pos.Y+10 > float32(dim.Size.Y) {
					ui.IsAdjustThe = UIDragTheRE & ui.IsAdjustMask
					pointer.CursorNorthWestResize.Add(gtx.Ops)
				} else if pos.X-10 < -5 {
					ui.IsAdjustThe = UIDragTheW & ui.IsAdjustMask
					pointer.CursorColResize.Add(gtx.Ops)
				} else if pos.X+10 > float32(dim.Size.X) {
					ui.IsAdjustThe = UIDragTheE & ui.IsAdjustMask
					pointer.CursorColResize.Add(gtx.Ops)
				} else if pos.Y+10 > float32(dim.Size.Y) {
					ui.IsAdjustThe = UIDragTheR & ui.IsAdjustMask
					pointer.CursorRowResize.Add(gtx.Ops)
				} else if pos.Y-10 < -5 {
					ui.IsAdjustThe = UIDragTheN & ui.IsAdjustMask
					pointer.CursorRowResize.Add(gtx.Ops)
				} else {
					ui.IsAdjustThe = UIDragTheDragging & ui.IsAdjustMask
					pointer.CursorGrabbing.Add(gtx.Ops)
				}
				// 处理事件
				if event.Buttons == pointer.ButtonPrimary && event.Kind == pointer.Press {
					switch ui.IsAdjustThe {
					case UIDragTheN:
						ui.dragging = true
						ui.position.Y = pos.Y
					case UIDragTheR:
						ui.dragging = true
						ui.position.Y = float32(dim.Size.Y) - pos.Y
					case UIDragTheE:
						ui.dragging = true
						ui.position.X = float32(dim.Size.X) - pos.X
					case UIDragTheW:
						ui.dragging = true
						ui.position.X = pos.X
					case UIDragTheRE:
						ui.dragging = true
						ui.position.Y = float32(dim.Size.Y) - pos.Y
						ui.position.X = float32(dim.Size.X) - pos.X
					case UIDragTheDragging:
						ui.dragging = true
						ui.position = pos
					}
				}
			} else {
				// 处理移动过程
				if event.Kind == pointer.Drag {
					r := pos.Sub(ui.position)
					switch ui.IsAdjustThe {
					case UIDragTheN:
						point.Y += int(r.Y)
						dim.Size.Y -= int(r.Y)
						pointer.CursorRowResize.Add(gtx.Ops)
					case UIDragTheR:
						dim.Size.Y = int(pos.Y - ui.position.Y)
						pointer.CursorRowResize.Add(gtx.Ops)
					case UIDragTheE:
						dim.Size.X = int(pos.X - ui.position.X)
						pointer.CursorColResize.Add(gtx.Ops)
					case UIDragTheW:
						point.X += int(r.X)
						dim.Size.X -= int(r.X)
						pointer.CursorColResize.Add(gtx.Ops)
					case UIDragTheRE:
						dim.Size.Y = int(pos.Y - ui.position.Y)
						dim.Size.X = int(pos.X - ui.position.X)
						pointer.CursorNorthWestResize.Add(gtx.Ops)
					case UIDragTheDragging:
						point.X += int(r.X)
						point.Y += int(r.Y)
						pointer.CursorGrabbing.Add(gtx.Ops)
					}
				} else {
					ui.Use = false
					ui.dragging = false
					ui.IsAdjustThe = 0
					ui.SetDisable(false)
					pointer.CursorDefault.Add(gtx.Ops)
				}
			}
		}
	}
	if ui.dragging {
		ui.LayoutFace.Update()

	}
	return ui.LayoutFace.Layout(gtx)
}
