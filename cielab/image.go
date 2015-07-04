package cielab

import (
	"image"
	"image/color"
)

// Color represents a color in LAB
type Color struct{ L, A, B float64 }

var ColorModel = color.ModelFunc(convert)

func convert(c color.Color) color.Color {
	if _, ok := c.(Color); ok {
		return c
	}

	r, g, b, _ := c.RGBA()
	l, la, lb := SRGB16ToLAB(r, g, b)
	return Color{l, la, lb}
}

func (c Color) RGBA() (r, g, b, a uint32) {
	r, g, b = LABToSRGB16(c.L, c.A, c.B)
	return r, g, b, 0xffff
}

// Image implements an image in CIELab color space
type Image struct {
	// L, A, B hold the appropriate channels
	L, A, B []float64
	// Stride is the L, A and B stride in pixels
	Stride int
	// Image bounds
	Rect image.Rectangle
}

func NewLAB(r image.Rectangle) *Image {
	w, h := r.Dx(), r.Dy()
	sz := w * h * 4
	return &Image{
		L: make([]float64, sz),
		A: make([]float64, sz),
		B: make([]float64, sz),

		Stride: w,
		Rect:   r,
	}
}

func (p *Image) Clone() *Image {
	m := &Image{}
	*m = *p

	m.L = make([]float64, len(p.L))
	copy(m.L, p.L)
	m.A = make([]float64, len(p.A))
	copy(m.A, p.A)
	m.B = make([]float64, len(p.B))
	copy(m.B, p.B)

	return m
}

func (p *Image) ColorModel() color.Model { return ColorModel }

func (p *Image) Bounds() image.Rectangle { return p.Rect }

func (p *Image) At(x, y int) color.Color { return p.LABAt(x, y) }

func (p *Image) LABAt(x, y int) Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return Color{}
	}
	i := p.Offset(x, y)
	return Color{p.L[i], p.A[i], p.B[i]}
}

func (p *Image) Offset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x - p.Rect.Min.X)
}

func (p *Image) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}

	i := p.Offset(x, y)
	c1 := ColorModel.Convert(c).(Color)
	p.L[i], p.A[i], p.B[i] = c1.L, c1.A, c1.B
}

func (p *Image) SetLAB(x, y int, c Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}

	i := p.Offset(x, y)
	p.L[i], p.A[i], p.B[i] = c.L, c.A, c.B
}
