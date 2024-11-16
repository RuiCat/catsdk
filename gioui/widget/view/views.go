package view

import (
	"gioui/anime/canvas"
	"gioui/io/event"
	"gioui/io/key"
	"gioui/io/pointer"
	"gioui/op"
	"gioui/op/clip"
	"gioui/op/paint"
	"gioui/widget"
	"gioui/widget/layout"
	"gioui/widget/material"
	"image"
	"mat/mat/spatial/f32"
)

// ViewsNew 创建控件
func ViewsNew(th *material.Theme, X, Y, Width, Height int, Size image.Point) *Views {
	views := &Views{
		Theme: th,
		Size:  Size,
		Position: image.Point{
			X: X,
			Y: Y,
		},
		Dimensions: layout.Dimensions{
			Baseline: 2,
		},
		UseGrid:      true,
		GridDistance: 50,
	}
	views.SetSize(Width, Height)
	views.affine.Scale = 1
	views.Scrollbar[0].Use = true
	views.Scrollbar[0].ScrollbarStyle = material.Scrollbar(views.Theme, &views.Scrollbar[0].Scrollbar)
	views.Scrollbar[1].Use = true
	views.Scrollbar[1].ScrollbarStyle = material.Scrollbar(views.Theme, &views.Scrollbar[1].Scrollbar)
	views.Background()
	return views
}

// Views 视图控件
type Views struct {
	*material.Theme                      // 主题
	Size               image.Point       // 组件内画板大小
	Position           image.Point       // 组件位置
	Dimensions         layout.Dimensions // 组件大小
	IsSizeCalculation  bool              // 计算大小
	UseGrid            bool              // 启用网格
	GridDistance       int               // 网格大小
	ComponentList      []LayoutFace      // 布局组件
	componentListIndex []int             // 组件索引
	Scrollbar          [2]struct {
		Use                     bool    // 滚动条生效
		Distance                float32 // 滚动条滚动距离
		widget.Scrollbar                // 滚动条对象
		material.ScrollbarStyle         // 滚动条对象风格
	}
	stroke   op.CallOp // 组件背景网格与边框
	clipRect clip.Op   // 用于组件剪裁
	moves    struct {
		Use      bool      // 背景移动状态标记
		Position f32.Point // 背景移动位置记录
	}
	affine struct {
		Scale        float32 // 缩放倍数
		f32.Affine2D         // 缩放矩阵
	}
}

// SetSize 设置大小
func (v *Views) SetSize(Width, Height int) {
	v.Dimensions.Size = image.Point{X: Width, Y: Height}
	v.clipRect = clip.Rect{Max: v.Dimensions.Size}.Op()
	v.Background()
}

// Add 添加组件
func (v *Views) Add(point image.Point, size image.Point, face UILayoutFace) {
	ui := face.GetUILayout()
	ui.Point = &image.Point{X: point.X, Y: point.Y}
	ui.Dimensions = &layout.Dimensions{Size: image.Point{X: size.X, Y: size.Y}}
	v.ComponentList = append(v.ComponentList, face)
	v.Scale(0)
}

// Scale 组件缩放
func (v *Views) Scale(scale float32) {
	if scale > 0 {
		v.affine.Scale = scale
	}
	v.affine.Affine2D = f32.Affine2D{}.Scale(
		f32.Point{X: 1, Y: 1},
		f32.Point{X: v.affine.Scale, Y: v.affine.Scale})
	v.Background()
}

// Layout 布局
func (v *Views) Layout(gtx layout.Context) layout.Dimensions {
	// 绘制边框
	paint.FillShape(gtx.Ops, v.Fg, clip.Stroke{Path: clip.RRect{
		Rect: image.Rectangle{
			Max: v.Dimensions.Size,
		},
	}.Path(&op.Ops{}), Width: 1.5}.Op())
	AffinePop := op.Affine(v.affine.Affine2D).Push(gtx.Ops) // 缩放剩余对象
	defer func() {
		AffinePop.Pop() // 还原组件缩放
		// 绘制滚动条
		gtx.Constraints.Max = v.Dimensions.Size.Sub(image.Point{X: 5, Y: 5}) // 定义长度并且减去右下角的空余位置偏移
		if v.Scrollbar[0].Use {
			TransformStack := op.Offset(image.Point{Y: v.Dimensions.Size.Y - 10}).Push(gtx.Ops)
			v.Scrollbar[0].ScrollbarStyle.Layout(gtx, layout.Horizontal, v.Scrollbar[0].Distance, v.Scrollbar[0].Distance)
			TransformStack.Pop()
		}
		if v.Scrollbar[1].Use {
			TransformStack := op.Offset(image.Point{X: v.Dimensions.Size.X - 10}).Push(gtx.Ops)
			v.Scrollbar[1].ScrollbarStyle.Layout(gtx, layout.Vertical, v.Scrollbar[1].Distance, v.Scrollbar[1].Distance)
			TransformStack.Pop()
		}
		v.clipRect.Push(gtx.Ops).Pop() // 剪裁组件位置
	}()
	// 判断组件大小是否扩展到 gtx 上下文大小
	size := v.Dimensions.Size
	if v.IsSizeCalculation {
		if size.X != gtx.Constraints.Max.X || size.Y != gtx.Constraints.Max.Y {
			v.SetSize(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
			return v.Dimensions
		}
		size = size.Sub(v.Position)
	}
	// 偏移位置
	offest := f32.Point{
		X: v.Scrollbar[0].Distance * float32(v.Size.X),
		Y: v.Scrollbar[1].Distance * float32(v.Size.Y),
	}
	// 处理事件
	event.Op(gtx.Ops, v)
	gtx.Execute(key.FocusCmd{Tag: v})
	for {
		if ev, ok := gtx.Event(
			pointer.Filter{
				Target:  v,
				Kinds:   pointer.Cancel | pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Enter | pointer.Leave | pointer.Scroll,
				ScrollY: pointer.ScrollRange{Max: 1},
			},
			key.FocusFilter{Target: v},                  // 接收按键
			key.Filter{Focus: v, Optional: key.ModCtrl}, // Ctrl 按下
		); ok {
			switch event := ev.(type) {
			case pointer.Event:
				pos := event.Position.Sub(f32.Point{X: float32(v.Position.X), Y: float32(v.Position.Y)})
				if pos.X > 0 && pos.Y > 0 && pos.X < float32(size.X)-10 && pos.Y < float32(size.Y)-10 || v.moves.Use {
					// 背景移动实现
					if event.Buttons == pointer.ButtonTertiary {
						switch event.Kind {
						case pointer.Press:
							if !v.moves.Use {
								v.moves.Position = pos.Add(offest)
							}
							v.moves.Use = true
						case pointer.Drag:
							if v.moves.Use {
								pos := pos.Sub(v.moves.Position)
								v.Scrollbar[0].Distance = -pos.X / float32(v.Size.X)
								v.Scrollbar[1].Distance = -pos.Y / float32(v.Size.Y)
							}
						}
					} else {
						v.moves.Use = false
					}
				}
				// 动态缩放
				if event.Kind == pointer.Scroll {
					if event.Scroll.Y == 1 {
						v.affine.Scale -= 0.01
					} else {
						v.affine.Scale += 0.01
					}
					v.Scale(mix(0.1, v.affine.Scale, 2))
				}
			case key.Event: // 键盘事件
			}
		} else {
			break
		}
	}
	// 更新显示组件
	v.componentListIndex = v.componentListIndex[:0]
	// 计算内容
	for i, cl := range v.ComponentList {
		// 计算显示矩阵
		size := cl.GetDimensions().Size
		position := cl.GetPoint()
		if position.X <= int(offest.X)+v.Dimensions.Size.X &&
			position.X+size.X >= int(offest.X) &&
			position.Y <= int(offest.Y)+v.Dimensions.Size.Y &&
			position.Y+size.Y >= int(offest.Y) {
			v.componentListIndex = append(v.componentListIndex, i)
		}
	}
	// 计算滚动条位置
	if !v.moves.Use {
		v.Scrollbar[1].Distance += v.Scrollbar[1].Scrollbar.ScrollDistance()
		v.Scrollbar[0].Distance += v.Scrollbar[0].Scrollbar.ScrollDistance()
	}
	v.Scrollbar[0].Distance = mix(0, v.Scrollbar[0].Distance, 1)
	v.Scrollbar[1].Distance = mix(0, v.Scrollbar[1].Distance, 1)
	// 绘制控件
	defer op.Offset(v.Position).Push(gtx.Ops).Pop() // 移动画布
	// 绘制背景网格
	if v.UseGrid {
		offest := op.Offset(image.Point{
			X: -int(v.Scrollbar[0].Distance*float32(v.Size.X)) % v.GridDistance,
			Y: -int(v.Scrollbar[1].Distance*float32(v.Size.Y)) % v.GridDistance,
		}).Push(gtx.Ops)
		v.stroke.Add(gtx.Ops)
		offest.Pop()
	}
	// 绘制组件
	OffsetPop := op.Offset(image.Point{
		X: -int(offest.X),
		Y: -int(offest.Y)},
	).Push(gtx.Ops) // 移动剪裁位置
	// 绘图边界
	paint.FillShape(gtx.Ops, v.ContrastBg, clip.Stroke{Path: clip.RRect{
		Rect: image.Rectangle{
			Max: v.Size,
		},
	}.Path(&op.Ops{}), Width: 2.5}.Op())
	for _, i := range v.componentListIndex {
		cl := v.ComponentList[i]
		offset := op.Offset(*cl.GetPoint()).Push(gtx.Ops)
		offsetRect := clip.Rect{Max: cl.GetDimensions().Size}.Push(gtx.Ops)
		cl.Layout(gtx)
		offsetRect.Pop()
		offset.Pop()
	}
	OffsetPop.Pop() // 还原组件偏移移动
	return v.Dimensions
}

// Background 更新背景
func (v *Views) Background() {
	if v.UseGrid {
		zy := v.GridDistance
		zx := float64(zy / 5)
		size := v.Dimensions.Size
		if v.affine.Scale < 1 {
			size = size.Mul(max(int(float32(v.GridDistance)/(v.affine.Scale)), v.GridDistance) / 10)
		}
		x, X := size.X/zy+2, size.X+zy
		y, Y := size.Y/zy+2, size.Y+zy
		gioGtx := &canvas.OpsGio{}
		ctx := gioGtx.Gio(image.Point{X, Y})
		ctx.ReflectY()
		ctx.Translate(0, -float64(Y))
		defer func() {
			ctx.FillStroke()
			v.stroke = gioGtx.MacroOp.Stop()
		}()
		ctx.SetFillColor(canvas.Transparent)
		// 绘制背景网格
		gridsX, gridsXZ := &canvas.Path{}, &canvas.Path{}
		for i := 0; i < x; i++ {
			ix := float64(i * zy)
			gridsX.MoveTo(ix, 0)
			gridsX.LineTo(ix, float64(Y))
			iz := float64((i + 1) * zy)
			for ; ix < iz; ix += zx {
				gridsXZ.MoveTo(ix, 0)
				gridsXZ.LineTo(ix, float64(Y))
			}
		}
		gridsY, gridsYZ := &canvas.Path{}, &canvas.Path{}
		for i := 0; i < y; i++ {
			iy := float64(i * zy)
			gridsY.MoveTo(0, iy)
			gridsY.LineTo(float64(X), iy)
			iz := float64((i + 1) * zy)
			for ; iy < iz; iy += zx {
				gridsYZ.MoveTo(0, iy)
				gridsYZ.LineTo(float64(X), iy)
			}
		}
		// 开始绘制
		ctx.SetStrokeColor(v.Fg)
		ctx.SetStrokeWidth(1.5)
		ctx.DrawPath(0, 0, gridsX)
		ctx.DrawPath(0, 0, gridsY)
		ctx.SetStrokeColor(v.ContrastBg)
		ctx.SetStrokeWidth(0.5)
		ctx.DrawPath(0, 0, gridsXZ)
		ctx.DrawPath(0, 0, gridsYZ)
	}
}

func mix(a0, x, a1 float32) float32 {
	if x <= a0 {
		return a0
	} else if x >= a1 {
		return a1
	}
	return x
}
