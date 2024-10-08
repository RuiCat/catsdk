// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vggio // import "gioui/anime/plot/vg/vggio"

import (
	"image/color"

	"gioui/op"

	"gioui/anime/plot/vg"
	"mat/mat/spatial/f32"
)

// ctxops holds a stack of Gio operations.
type ctxops struct {
	ops *op.Ops   // ops is the Gio operations vggio is drawing on.
	ctx []context // ctx is the stack of Gio operations vggio is manipulating.

	w   vg.Length // w is the canvas window width.
	h   vg.Length // h is the canvas window height.
	dpi float64   // dpi is the canvas window dots per inch resolution.
}

// context holds state about the Gio backing store.
// context provides methods to translate between Gio values (reference frame,
// operations and stack) and their plot/vg counterparts.
type context struct {
	color   color.Color // color is the current color.
	linew   vg.Length   // linew is the current line width.
	pattern []vg.Length // pattern is the current line style.
	offset  vg.Length   // offset is the current line style.

	trans op.TransformStack // trans is the Gio transform context stack.
}

func (ctx *ctxops) cur() *context {
	return &ctx.ctx[len(ctx.ctx)-1]
}

func (ctx *ctxops) push() {
	ctx.ctx = append(ctx.ctx, *ctx.cur())
	ctx.cur().trans = op.TransformOp{}.Push(ctx.ops)
}

func (ctx *ctxops) pop() {
	ctx.cur().trans.Pop()
	ctx.ctx = ctx.ctx[:len(ctx.ctx)-1]
}

func (ctx *ctxops) scale(x, y float64) {
	op.Affine(f32.Affine2D{}.Scale(
		f32.Pt(0, 0),
		f32.Pt(float32(x), float32(y)),
	)).Add(ctx.ops)
}

func (ctx *ctxops) translate(x, y float64) {
	op.Affine(f32.Affine2D{}.Offset(
		f32.Pt(float32(x), float32(y)),
	)).Add(ctx.ops)
}

func (ctx *ctxops) rotate(rad float64) {
	op.Affine(f32.Affine2D{}.Rotate(
		f32.Pt(0, 0), float32(rad),
	)).Add(ctx.ops)
}

func (ctx *ctxops) invertY() {
	ctx.translate(0, ctx.h.Dots(ctx.dpi))
	ctx.scale(1, -1)
}

func (ctx *ctxops) pt32(p vg.Point) f32.Point {
	return f32.Point{
		X: float32(p.X.Dots(ctx.dpi)),
		Y: float32(p.Y.Dots(ctx.dpi)),
	}
}
