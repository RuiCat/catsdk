package draw

import (
	"image/color"
	"mat/asm/f32"
	"math"
)

var (
	Black = Color{0, 0, 0, 0}
	White = Color{1, 1, 1, 1}
)

type Color struct {
	R, G, B, A float32
}

func HexColor(x int) Color {
	r := float32((x>>16)&0xff) / 255
	g := float32((x>>8)&0xff) / 255
	b := float32((x>>0)&0xff) / 255
	return Color{r, g, b, 0.5}.Pow(2.2)
}

func Kelvin(K float32) Color {
	k := float64(K)
	var red, green, blue float64
	// red
	if k >= 6600 {
		a := 351.97690566805693
		b := 0.114206453784165
		c := -40.25366309332127
		x := k/100 - 55
		red = a + b*x + c*math.Log(x)
	} else {
		red = 255
	}
	// green
	if k >= 6600 {
		a := 325.4494125711974
		b := 0.07943456536662342
		c := -28.0852963507957
		x := k/100 - 50
		green = a + b*x + c*math.Log(x)
	} else if k >= 1000 {
		a := -155.25485562709179
		b := -0.44596950469579133
		c := 104.49216199393888
		x := k/100 - 2
		green = a + b*x + c*math.Log(x)
	} else {
		green = 0
	}
	// blue
	if k >= 6600 {
		blue = 255
	} else if k >= 2000 {
		a := -254.76935184120902
		b := 0.8274096064007395
		c := 115.67994401066147
		x := k/100 - 10
		blue = a + b*x + c*math.Log(x)
	} else {
		blue = 0
	}
	red = math.Min(1, red/255)
	green = math.Min(1, green/255)
	blue = math.Min(1, blue/255)
	return Color{float32(red), float32(green), float32(blue), 1}
}

func NewColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	return Color{float32(r) / 65535, float32(g) / 65535, float32(b) / 65535, float32(a) / 65535}
}

func (a Color) RGBA() color.RGBA {
	r := uint8(math.Max(0, math.Min(255, float64(a.R)*255)))
	g := uint8(math.Max(0, math.Min(255, float64(a.G)*255)))
	b := uint8(math.Max(0, math.Min(255, float64(a.B)*255)))
	ar := uint8(math.Max(0, math.Min(255, float64(a.A)*255)))
	return color.RGBA{r, g, b, ar}
}

func (a Color) RGBA64() color.RGBA64 {
	r := uint16(math.Max(0, math.Min(65535, float64(a.R)*65535)))
	g := uint16(math.Max(0, math.Min(65535, float64(a.G)*65535)))
	b := uint16(math.Max(0, math.Min(65535, float64(a.B)*65535)))
	ar := uint16(math.Max(0, math.Min(65535, float64(a.A)*65535)))
	return color.RGBA64{r, g, b, ar}
}

func (a Color) Add(b Color) Color {
	return Color{a.R + b.R, a.G + b.G, a.B + b.B, a.A + b.A}
}

func (a Color) Sub(b Color) Color {
	return Color{a.R - b.R, a.G - b.G, a.B - b.B, a.A - b.A}
}

func (a Color) Mul(b Color) Color {
	return Color{a.R * b.R, a.G * b.G, a.B * b.B, a.A * b.A}
}

func (a Color) MulScalar(b float32) Color {
	return Color{a.R * b, a.G * b, a.B * b, a.A * b}
}

func (a Color) DivScalar(b float32) Color {
	return Color{a.R / b, a.G / b, a.B / b, a.A / b}
}

func (a Color) Min(b Color) Color {
	return Color{f32.Min(a.R, b.R), f32.Min(a.G, b.G), f32.Min(a.B, b.B), f32.Max(a.A, b.A)}
}

func (a Color) Max(b Color) Color {
	return Color{f32.Max(a.R, b.R), f32.Max(a.G, b.G), f32.Max(a.B, b.B), f32.Max(a.A, b.A)}
}

func (a Color) MinComponent() float32 {
	return f32.Min(f32.Min(a.R, a.G), a.B)
}

func (a Color) MaxComponent() float32 {
	return f32.Max(f32.Max(a.R, a.G), a.B)
}

func (a Color) Pow(b float32) Color {
	return Color{f32.Pow(a.R, b), f32.Pow(a.G, b), f32.Pow(a.B, b), f32.Pow(a.A, b)}
}

func (a Color) Mix(b Color, pct float32) Color {
	a = a.MulScalar(1 - pct)
	b = b.MulScalar(pct)
	return a.Add(b)
}

func (a Color) AddColor(b Color) Color {
	aR, aG, aB, aA := a.R, a.G, a.B, a.A
	bR, bG, bB, bA := b.R, b.G, b.B, b.A
	R := aR*aA + bR*bA*(1-aA)
	G := aG*aA + bG*bA*(1-aA)
	B := aB*aA + bB*bA*(1-aA)
	Alpha := 1 - (1-aA)*(1-bA)
	return Color{R / Alpha, G / Alpha, B / Alpha, Alpha}
}
