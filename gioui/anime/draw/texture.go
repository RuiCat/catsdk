package draw

import (
	"image"
	"mat/asm/f32"
)

type Texture interface {
	Sample(u, v float64) Color
	NormalSample(u, v float64) Vector
	BumpSample(u, v float64) Vector
	Pow(a float32) Texture
	MulScalar(a float32) Texture
}

var textures map[string]Texture

func init() {
	textures = make(map[string]Texture)
}

func GetTexture(path string) Texture {
	if texture, ok := textures[path]; ok {
		return texture
	}
	if texture, err := LoadTexture(path); err == nil {
		textures[path] = texture
		return texture
	}
	return nil
}

func LoadTexture(path string) (Texture, error) {
	im, err := LoadImage(path)
	if err != nil {
		return nil, err
	}
	return NewTexture(im), nil
}

type ColorTexture struct {
	Width  int
	Height int
	Data   []Color
}

func NewTexture(im image.Image) Texture {
	size := im.Bounds().Max
	data := make([]Color, size.X*size.Y)
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			index := y*size.X + x
			data[index] = NewColor(im.At(x, y)).Pow(2.2)
		}
	}
	return &ColorTexture{size.X, size.Y, data}
}

func (t *ColorTexture) Pow(a float32) Texture {
	for i := range t.Data {
		t.Data[i] = t.Data[i].Pow(a)
	}
	return t
}

func (t *ColorTexture) MulScalar(a float32) Texture {
	for i := range t.Data {
		t.Data[i] = t.Data[i].MulScalar(a)
	}
	return t
}

func (t *ColorTexture) bilinearSample(u, v float32) Color {
	if u == 1 {
		u -= EPS
	}
	if v == 1 {
		v -= EPS
	}
	w := float32(t.Width) - 2
	h := float32(t.Height) - 2
	X, x := f32.Modf(u * w)
	Y, y := f32.Modf(v * h)
	x0 := int(X)
	y0 := int(Y)
	x1 := x0 + 1
	y1 := y0 + 1
	c00 := t.Data[y0*t.Width+x0]
	c01 := t.Data[y1*t.Width+x0]
	c10 := t.Data[y0*t.Width+x1]
	c11 := t.Data[y1*t.Width+x1]
	c := Black
	c = c.Add(c00.MulScalar((1 - x) * (1 - y)))
	c = c.Add(c10.MulScalar(x * (1 - y)))
	c = c.Add(c01.MulScalar((1 - x) * y))
	c = c.Add(c11.MulScalar(x * y))
	return c
}

func (t *ColorTexture) Sample(u, v float64) Color {
	return t.bilinearSample(Fract(Fract(float32(u))+1), 1-Fract(Fract(float32(v))+1))
}

func (t *ColorTexture) NormalSample(u, v float64) Vector {
	c := t.Sample(u, v)
	return Vector{X: float64(c.R)*2 - 1, Y: float64(c.G)*2 - 1, Z: float64(c.B)*2 - 1}.Normalize()
}

func (t *ColorTexture) BumpSample(u, v float64) Vector {
	u = float64(Fract(Fract(float32(u)) + 1))
	v = float64(Fract(Fract(float32(v)) + 1))
	v = 1 - v
	x := int(u * float64(t.Width))
	y := int(v * float64(t.Height))
	x1, x2 := ClampInt(x-1, 0, t.Width-1), ClampInt(x+1, 0, t.Width-1)
	y1, y2 := ClampInt(y-1, 0, t.Height-1), ClampInt(y+1, 0, t.Height-1)
	cx := t.Data[y*t.Width+x1].Sub(t.Data[y*t.Width+x2])
	cy := t.Data[y1*t.Width+x].Sub(t.Data[y2*t.Width+x])
	return Vector{X: float64(cx.R), Y: float64(cy.R), Z: 0}
}

type ColorImage struct {
	Width  int
	Height int
	image.Image
}

func NewColorImage(im image.Image) *ColorImage {
	size := im.Bounds().Max
	return &ColorImage{size.X, size.Y, im}
}

func (t *ColorImage) bilinearSample(u, v float32) Color {
	if u == 1 {
		u -= EPS
	}
	if v == 1 {
		v -= EPS
	}
	w := float32(t.Width) - 2
	h := float32(t.Height) - 2
	X, _ := f32.Modf(u * w)
	Y, _ := f32.Modf(v * h)
	x0 := int(X)
	y0 := int(Y)
	return NewColor(t.At(x0, y0)).Pow(2.2)
}
func (t *ColorImage) Sample(u, v float64) Color {
	return t.bilinearSample(Fract(Fract(float32(u))+1), 1-Fract(Fract(float32(v))+1))
}
func (t *ColorImage) NormalSample(u, v float64) Vector { return Vector{} }
func (t *ColorImage) BumpSample(u, v float64) Vector   { return Vector{} }
func (t *ColorImage) Pow(a float32) Texture            { return t }
func (t *ColorImage) MulScalar(a float32) Texture      { return t }
