package draw

import (
	"image"
	"mat/asm/f32"
)

// Channel 是一个用于表示图形缓存通道的类型
type Channel int

// ColorChannel 表示颜色通道
const (
	ColorChannel             Channel = iota // 颜色通道的编号
	VarianceChannel                         // 方差通道的编号
	StandardDeviationChannel                // 标准差通道的编号
	SamplesChannel                          // 样本通道的编号
)

type Pixel struct {
	Samples int
	M, V    Color
}

func (p *Pixel) AddSample(sample Color) {
	p.Samples++
	if p.Samples == 1 {
		p.M = sample
		return
	}
	m := p.M
	p.M = p.M.Add(sample.Sub(p.M).DivScalar(float32(p.Samples)))
	p.V = p.V.Add(sample.Sub(m).Mul(sample.Sub(p.M)))
}

func (p *Pixel) Color() Color {
	return p.M
}

func (p *Pixel) Variance() Color {
	if p.Samples < 2 {
		return Black
	}
	return p.V.DivScalar(float32(p.Samples - 1))
}

func (p *Pixel) StandardDeviation() Color {
	return p.Variance().Pow(0.5)
}

type Buffer struct {
	W, H   int
	Pixels []Pixel
}

func NewBuffer(w, h int) *Buffer {
	pixels := make([]Pixel, w*h)
	return &Buffer{w, h, pixels}
}

func (b *Buffer) Copy() *Buffer {
	pixels := make([]Pixel, b.W*b.H)
	copy(pixels, b.Pixels)
	return &Buffer{b.W, b.H, pixels}
}

func (b *Buffer) SetLen(w, h int) {
	le := w * h
	if len(b.Pixels) < le {
		b.Pixels = append(b.Pixels, make([]Pixel, le-len(b.Pixels))...)
	}
	clear(b.Pixels)
	b.W, b.H = w, h
}

func (b *Buffer) AddSample(x, y int, sample Color) {
	b.Pixels[y*b.W+x].AddSample(sample)
}

func (b *Buffer) Samples(x, y int) int {
	return b.Pixels[y*b.W+x].Samples
}

func (b *Buffer) Color(x, y int) Color {
	return b.Pixels[y*b.W+x].Color()
}

func (b *Buffer) Variance(x, y int) Color {
	return b.Pixels[y*b.W+x].Variance()
}

func (b *Buffer) StandardDeviation(x, y int) Color {
	return b.Pixels[y*b.W+x].StandardDeviation()
}

func (b *Buffer) Image(channel Channel) image.Image {
	result := image.NewRGBA64(image.Rect(0, 0, b.W, b.H))
	var maxSamples float32
	if channel == SamplesChannel {
		for _, pixel := range b.Pixels {
			maxSamples = f32.Max(maxSamples, float32(pixel.Samples))
		}
	}
	for y := 0; y < b.H; y++ {
		for x := 0; x < b.W; x++ {
			var c Color
			switch channel {
			case ColorChannel:
				c = b.Pixels[y*b.W+x].Color().Pow(1 / 2.2)
			case VarianceChannel:
				c = b.Pixels[y*b.W+x].Variance()
			case StandardDeviationChannel:
				c = b.Pixels[y*b.W+x].StandardDeviation()
			case SamplesChannel:
				p := float32(b.Pixels[y*b.W+x].Samples) / maxSamples
				c = Color{p, p, p, p}
			}
			result.SetRGBA64(x, y, c.RGBA64())
		}
	}
	return result
}
