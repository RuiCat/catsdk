package drawing

import (
	"image"
	"image/color"
	"sdk/freetype/raster"
	"sdk/x/image/draw"
	"sdk/x/image/math/f64"
)

// Painter implements the freetype raster.Painter and has a SetColor method like the RGBAPainter
type Painter interface {
	raster.Painter
	SetColor(color color.Color)
}

// DrawImage draws an image into dest using an affine transformation matrix, an op and a filter
func DrawImage(src image.Image, dest draw.Image, tr Matrix, op draw.Op, filter ImageFilter) {
	var transformer draw.Transformer
	switch filter {
	case LinearFilter:
		transformer = draw.NearestNeighbor
	case BilinearFilter:
		transformer = draw.BiLinear
	case BicubicFilter:
		transformer = draw.CatmullRom
	}
	transformer.Transform(dest, f64.Aff3{tr[0], tr[1], tr[4], tr[2], tr[3], tr[5]}, src, src.Bounds(), op, nil)
}
