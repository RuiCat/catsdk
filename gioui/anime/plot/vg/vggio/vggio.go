// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vggio provides a vg.Canvas implementation backed by Gio,
// a toolkit that implements portable immediate GUI mode in Go.
//
// More informations about Gio can be found at https://gioui.org/.
package vggio // import "gioui/anime/plot/vg/vggio"

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"strings"
	"sync"

	giofont "gioui/font"
	"gioui/font/opentype"
	"gioui/gpu/headless"

	"gioui/op"
	"gioui/op/clip"
	"gioui/op/paint"
	"gioui/widget/layout"
	"gioui/widget/material"
	"gioui/widget/stroke"
	"gioui/widget/text"

	"sdk/x/image/draw"
	"sdk/x/image/font/sfnt"

	bstroke "sdk/stroke"

	"gioui/anime/plot/font"
	"gioui/anime/plot/vg"
	"mat/mat/spatial/f32"
	"mat/mat/spatial/unit"
)

var (
	_ vg.Canvas      = (*Canvas)(nil)
	_ vg.CanvasSizer = (*Canvas)(nil)
)

// Canvas implements the vg.Canvas interface,
// drawing to an image.Image using vgimg and painting that image
// into a Gioui context.
type Canvas struct {
	gtx layout.Context
	ctx ctxops

	bkg color.Color // bkg is the background color.
}

// DefaultDPI is the default dot resolution for image
// drawing in dots per inch.
const DefaultDPI = 96

// New returns a new image canvas with the provided dimensions and options.
// The currently accepted options are UseDPI and UseBackgroundColor.
// If the resolution or background color are not specified, defaults are used.
func New(gtx layout.Context, w, h vg.Length, opts ...option) *Canvas {
	cfg := &config{
		dpi: DefaultDPI,
		bkg: color.White,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	c := &Canvas{
		gtx: gtx,
		ctx: ctxops{
			ops: gtx.Ops,
			ctx: []context{
				{color: color.Black},
			},
			w:   w,
			h:   h,
			dpi: cfg.dpi,
		},
		bkg: cfg.bkg,
	}

	// flip the Y-axis so that Y grows from bottom to top and
	// Y=0 is at the bottom of the image.
	c.ctx.invertY()

	vg.Initialize(c)

	return c
}

type config struct {
	dpi float64
	bkg color.Color
}

type option func(*config)

// UseDPI sets the dots per inch of a canvas. It should only be
// used as an option argument when initializing a new canvas.
func UseDPI(dpi int) option {
	if dpi <= 0 {
		panic("DPI must be > 0.")
	}
	return func(c *config) {
		c.dpi = float64(dpi)
	}
}

// UseBackgroundColor specifies the image background color.
// Without UseBackgroundColor, the default color is white.
func UseBackgroundColor(c color.Color) option {
	return func(cfg *config) {
		cfg.bkg = c
	}
}

// Size implement vg.CanvasSizer.
func (c *Canvas) Size() (w, h vg.Length) {
	return c.ctx.w, c.ctx.h
}

// DPI returns the resolution of the receiver in pixels per inch.
func (c *Canvas) DPI() float64 {
	return c.ctx.dpi
}

// Paint returns the painting operations.
func (c *Canvas) Paint() *op.Ops {
	return c.gtx.Ops
}

// Screenshot returns a screenshot of the canvas as an image.
func (c *Canvas) Screenshot() (image.Image, error) {
	win, err := headless.NewWindow(
		int(c.ctx.w.Dots(c.ctx.dpi)),
		int(c.ctx.h.Dots(c.ctx.dpi)),
	)
	if err != nil {
		return nil, fmt.Errorf("vggio: could not create headless window: %w", err)
	}

	err = win.Frame(c.gtx.Ops)
	if err != nil {
		return nil, fmt.Errorf("vggio: could not run headless frame: %w", err)
	}

	img := image.NewRGBA(image.Rectangle{Max: win.Size()})
	err = win.Screenshot(img)
	if err != nil {
		return nil, fmt.Errorf("vggio: could not create screenshot: %w", err)
	}

	return img, nil
}

// SetLineWidth sets the width of stroked paths.
// If the width is not positive then stroked lines
// are not drawn.
//
// The initial line width is 1 point.
func (c *Canvas) SetLineWidth(w vg.Length) {
	c.ctx.cur().linew = w
}

// SetLineDash sets the dash pattern for lines.
// The pattern slice specifies the lengths of
// alternating dashes and gaps, and the offset
// specifies the distance into the dash pattern
// to start the dash.
//
// The initial dash pattern is a solid line.
func (c *Canvas) SetLineDash(pattern []vg.Length, offset vg.Length) {
	cur := c.ctx.cur()
	cur.pattern = pattern
	cur.offset = offset
}

// SetColor sets the current drawing color.
// Note that fill color and stroke color are
// the same, so if you want different fill
// and stroke colors then you must set a color,
// draw fills, set a new color and then draw lines.
//
// The initial color is black.
// If SetColor is called with a nil color then black is used.
func (c *Canvas) SetColor(clr color.Color) {
	if clr == nil {
		clr = color.Black
	}
	c.ctx.cur().color = clr
}

// Rotate applies a rotation transform to the context.
// The parameter is specified in radians.
func (c *Canvas) Rotate(rad float64) {
	c.ctx.rotate(rad)
}

// Translate applies a translational transform
// to the context.
func (c *Canvas) Translate(pt vg.Point) {
	c.ctx.translate(pt.X.Dots(c.ctx.dpi), pt.Y.Dots(c.ctx.dpi))
}

// Scale applies a scaling transform to the
// context.
func (c *Canvas) Scale(x, y float64) {
	c.ctx.scale(x, y)
}

// Push saves the current line width, the
// current dash pattern, the current
// transforms, and the current color
// onto a stack so that the state can later
// be restored by calling Pop().
func (c *Canvas) Push() {
	c.ctx.push()
}

// Pop restores the context saved by the
// corresponding call to Push().
func (c *Canvas) Pop() {
	c.ctx.pop()
}

// Stroke strokes the given path.
func (c *Canvas) Stroke(p vg.Path) {
	if c.ctx.cur().linew <= 0 {
		return
	}
	c.ctx.push()
	defer c.ctx.pop()

	var (
		cur    = c.ctx.cur()
		dashes stroke.Dashes
	)
	dashes.Phase = float32(cur.offset.Dots(c.ctx.dpi))
	dashes.Dashes = make([]float32, len(cur.pattern))
	for i, v := range cur.pattern {
		dashes.Dashes[i] = float32(v.Dots(c.ctx.dpi))
	}

	shape := stroke.Stroke{
		Path:   c.stroke(p),
		Width:  float32(cur.linew.Dots(c.ctx.dpi)),
		Cap:    stroke.FlatCap,
		Dashes: dashes,
	}.Op(c.ctx.ops)

	clr := c.ctx.cur().color
	paint.FillShape(c.ctx.ops, rgba(clr), shape)
}

// Fill fills the given path.
func (c *Canvas) Fill(p vg.Path) {
	c.ctx.push()
	defer c.ctx.pop()

	shape := clip.Outline{
		Path: c.outline(p),
	}.Op()

	clr := c.ctx.cur().color
	paint.FillShape(c.ctx.ops, rgba(clr), shape)
}

func rgba(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

func (c *Canvas) outline(p vg.Path) clip.PathSpec {
	var path clip.Path
	path.Begin(c.ctx.ops)
	for _, comp := range p {
		switch comp.Type {
		case vg.MoveComp:
			pt := c.ctx.pt32(comp.Pos)
			path.MoveTo(pt)

		case vg.LineComp:
			pt := c.ctx.pt32(comp.Pos)
			path.LineTo(pt)

		case vg.ArcComp:
			center := c.ctx.pt32(comp.Pos)
			path.ArcTo(center, center, float32(comp.Angle))

		case vg.CurveComp:
			switch len(comp.Control) {
			case 1:
				ctl := c.ctx.pt32(comp.Control[0])
				end := c.ctx.pt32(comp.Pos)
				path.QuadTo(ctl, end)
			case 2:
				ctl0 := c.ctx.pt32(comp.Control[0])
				ctl1 := c.ctx.pt32(comp.Control[1])
				end := c.ctx.pt32(comp.Pos)
				path.CubeTo(ctl0, ctl1, end)
			default:
				panic("vggio: invalid number of control points")
			}

		case vg.CloseComp:
			path.Close()

		default:
			panic(fmt.Sprintf("vggio: unknown path component %d", comp.Type))
		}
	}
	return path.End()
}

func (c *Canvas) stroke(p vg.Path) stroke.Path {
	var (
		path stroke.Path
		add  = func(seg stroke.Segment) {
			path.Segments = append(path.Segments, seg)
		}
		pen f32.Point
		beg f32.Point
	)

	for i, comp := range p {
		if i == 0 {
			beg = c.ctx.pt32(comp.Pos)
		}
		switch comp.Type {
		case vg.MoveComp:
			pt := c.ctx.pt32(comp.Pos)
			add(stroke.MoveTo(pt))
			pen = pt

		case vg.LineComp:
			pt := c.ctx.pt32(comp.Pos)
			add(stroke.LineTo(pt))
			pen = pt

		case vg.ArcComp:
			center := c.ctx.pt32(comp.Pos)
			arcs := arcTo(pen, center, center, float32(comp.Angle))
			path.Segments = append(path.Segments, xStroke(arcs)...)
			pen = f32.Point(arcs[len(arcs)-1].End)

		case vg.CurveComp:
			switch len(comp.Control) {
			case 1:
				var (
					ctl = c.ctx.pt32(comp.Control[0])
					end = c.ctx.pt32(comp.Pos)
				)
				add(stroke.QuadTo(ctl, end))
				pen = end
			case 2:
				var (
					ctl0 = c.ctx.pt32(comp.Control[0])
					ctl1 = c.ctx.pt32(comp.Control[1])
					end  = c.ctx.pt32(comp.Pos)
				)
				add(stroke.CubeTo(ctl0, ctl1, end))
				pen = end
			default:
				panic("vggio: invalid number of control points")
			}

		case vg.CloseComp:
			add(stroke.LineTo(beg))
			pen = beg

		default:
			panic(fmt.Sprintf("vggio: unknown path component %d", comp.Type))
		}
	}
	return path
}

// FillString fills in text at the specified
// location using the given font.
// If the font size is zero, the text is not drawn.
func (c *Canvas) FillString(fnt font.Face, pt vg.Point, txt string) {
	if fnt.Font.Size == 0 {
		return
	}
	c.ctx.push()
	defer c.ctx.pop()

	e := fnt.Extents()
	x := pt.X.Dots(c.ctx.dpi)
	y := pt.Y.Dots(c.ctx.dpi) - e.Descent.Dots(c.ctx.dpi)
	h := c.ctx.h.Dots(c.ctx.dpi)

	c.ctx.invertY()
	c.ctx.translate(x, h-y-fnt.Font.Size.Dots(c.ctx.dpi))

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(collectionFor(fnt)))
	lbl := material.Label(
		th,
		unit.Sp(float32(fnt.Font.Size.Dots(c.ctx.dpi))),
		txt,
	)
	lbl.Color = rgba(c.ctx.cur().color)
	lbl.Alignment = text.Start
	lbl.Layout(c.gtx)
}

// DrawImage draws the image, scaled to fit
// the destination rectangle.
func (c *Canvas) DrawImage(rect vg.Rectangle, img image.Image) {
	c.ctx.push()
	defer c.ctx.pop()

	var (
		ops    = c.ctx.ops
		dpi    = c.DPI()
		min    = rect.Min
		xmin   = min.X.Dots(dpi)
		ymin   = min.Y.Dots(dpi)
		rsz    = rect.Size()
		width  = rsz.X.Dots(dpi)
		height = rsz.Y.Dots(dpi)
		dst    = image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	)

	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Src, nil)

	c.ctx.scale(1, -1)
	c.ctx.translate(xmin, -ymin-height)
	paint.NewImageOp(dst).Add(ops)
	paint.PaintOp{}.Add(ops)
}

var dbfonts = &gioFontsCache{
	cache: make(map[string][]giofont.FontFace),
	fonts: make(map[string]struct{}),
}

type gioFontsCache struct {
	sync.RWMutex
	cache map[string][]giofont.FontFace
	fonts map[string]struct{}
	buf   sfnt.Buffer
}

func (cache *gioFontsCache) get(fnt font.Face) ([]giofont.FontFace, bool) {
	cache.RLock()
	defer cache.RUnlock()

	_, ok := cache.fonts[fnt.Name()]
	if !ok {
		return nil, false
	}
	name := collectionName(fnt.Name())
	return cache.cache[name], ok
}

func (cache *gioFontsCache) add(fnt font.Face) []giofont.FontFace {
	cache.Lock()
	defer cache.Unlock()

	name := fnt.Name()
	if fnt.Face == nil {
		panic(fmt.Errorf("vggio: nil plot/font.Face %q", name))
	}
	buf := new(bytes.Buffer)
	_, err := fnt.Face.WriteSourceTo(&cache.buf, buf)
	if err != nil {
		panic(fmt.Errorf("vggio: could not load font %q: %+v", name, err))
	}

	gioFace, err := opentype.Parse(buf.Bytes())
	if err != nil {
		panic(fmt.Errorf("vggio: could not parse font %q: %+v", name, err))
	}

	gioFnt := gonumToGioFont(fnt.Font)

	colName := collectionName(fnt.Name())
	cache.cache[colName] = append(cache.cache[colName], giofont.FontFace{
		Font: gioFnt,
		Face: gioFace,
	})
	cache.fonts[name] = struct{}{}

	return cache.cache[colName]
}

func gonumToGioFont(fnt font.Font) giofont.Font {
	o := giofont.Font{
		Typeface: giofont.Typeface(fnt.Typeface),
		Style:    giofont.Style(fnt.Style),
		Weight:   giofont.Weight(fnt.Weight),
	}
	return o
}

func collectionFor(fnt font.Face) []giofont.FontFace {
	coll, ok := dbfonts.get(fnt)
	if !ok {
		coll = dbfonts.add(fnt)
	}
	return coll
}

func collectionName(name string) string {
	// regroup fonts with name "Liberation-Italic", "Liberation-Bold", ...
	// under the same collection "Liberation".
	if strings.Contains(name, "-") {
		i := strings.Index(name, "-")
		name = name[:i]
	}
	return name
}

func arcTo(start, f1, f2 f32.Point, angle float32) []bstroke.Segment {
	if f1 == f2 {
		return bstroke.AppendArc(nil, bstroke.Pt(start.X, start.Y), bstroke.Pt(f1.X, f1.Y), angle)
	}
	return bstroke.AppendEllipticalArc(nil, bstroke.Pt(start.X, start.Y), bstroke.Pt(f1.X, f1.Y), bstroke.Pt(f2.X, f2.Y), angle)
}

func xStroke(bs []bstroke.Segment) []stroke.Segment {
	vs := make([]stroke.Segment, len(bs))
	for i, b := range bs {
		vs[i] = stroke.CubeTo(f32.Point(b.CP1), f32.Point(b.CP2), f32.Point(b.End))
	}
	return vs
}
