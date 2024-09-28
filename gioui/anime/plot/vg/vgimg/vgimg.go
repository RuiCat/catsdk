// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vgimg implements the vg.Canvas interface using
// git.sr.ht/~sbinet/gg as a backend to output raster images.
package vgimg // import "mat/plot/vg/vgimg"

import (
	"bufio"
	"gioui/anime/canvas"
	"gioui/anime/canvas/renderers"
	"gioui/anime/canvas/renderers/rasterizer"
	"gioui/anime/plot/vg"
	vgdraw "gioui/anime/plot/vg/draw"
	"image/jpeg"
	"image/png"
	"io"
	"sdk/x/image/tiff"
)

func init() {
	vgdraw.RegisterFormat("png", func(w, h vg.Length) vg.CanvasWriterTo {
		return PngCanvas{Canvas: New(w, h)}
	})

	vgdraw.RegisterFormat("jpg", func(w, h vg.Length) vg.CanvasWriterTo {
		return JpegCanvas{Canvas: New(w, h)}
	})

	vgdraw.RegisterFormat("jpeg", func(w, h vg.Length) vg.CanvasWriterTo {
		return JpegCanvas{Canvas: New(w, h)}
	})

	vgdraw.RegisterFormat("tif", func(w, h vg.Length) vg.CanvasWriterTo {
		return TiffCanvas{Canvas: New(w, h)}
	})

	vgdraw.RegisterFormat("tiff", func(w, h vg.Length) vg.CanvasWriterTo {
		return TiffCanvas{Canvas: New(w, h)}
	})
}

type Canvas struct {
	*renderers.GonumPlot
	Canvas *canvas.Canvas
}

func New(width, height vg.Length) *Canvas {
	c := &Canvas{}
	c.Canvas = canvas.New(float64(width), float64(height))
	c.GonumPlot = (renderers.NewGonumPlot(c.Canvas).Canvas).(*renderers.GonumPlot)
	return c
}

// WriterCounter implements the io.Writer interface, and counts
// the total number of bytes written.
type writerCounter struct {
	io.Writer
	n int64
}

func (w *writerCounter) Write(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	w.n += int64(n)
	return n, err
}

// A JpegCanvas is an image canvas with a WriteTo method
// that writes a jpeg image.
type JpegCanvas struct {
	*Canvas
}

// WriteTo implements the io.WriterTo interface, writing a jpeg image.
func (c JpegCanvas) WriteTo(w io.Writer) (int64, error) {
	wc := writerCounter{Writer: w}
	b := bufio.NewWriter(&wc)
	if err := jpeg.Encode(b, rasterizer.Draw(c.Canvas.Canvas, canvas.DefaultResolution, canvas.DefaultColorSpace), nil); err != nil {
		return wc.n, err
	}
	err := b.Flush()
	return wc.n, err
}

// A PngCanvas is an image canvas with a WriteTo method that
// writes a png image.
type PngCanvas struct {
	*Canvas
}

// WriteTo implements the io.WriterTo interface, writing a png image.
func (c PngCanvas) WriteTo(w io.Writer) (int64, error) {
	wc := writerCounter{Writer: w}
	b := bufio.NewWriter(&wc)
	if err := png.Encode(b, rasterizer.Draw(c.Canvas.Canvas, canvas.DefaultResolution, canvas.DefaultColorSpace)); err != nil {
		return wc.n, err
	}
	err := b.Flush()
	return wc.n, err
}

// A TiffCanvas is an image canvas with a WriteTo method that
// writes a tiff image.
type TiffCanvas struct {
	*Canvas
}

// WriteTo implements the io.WriterTo interface, writing a tiff image.
func (c TiffCanvas) WriteTo(w io.Writer) (int64, error) {
	wc := writerCounter{Writer: w}
	b := bufio.NewWriter(&wc)
	if err := tiff.Encode(b, rasterizer.Draw(c.Canvas.Canvas, canvas.DefaultResolution, canvas.DefaultColorSpace), nil); err != nil {
		return wc.n, err
	}
	err := b.Flush()
	return wc.n, err
}
