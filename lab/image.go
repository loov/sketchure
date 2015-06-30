package lab

import (
	"image"
	"image/color"
)

var Model = color.ModelFunc(labmodel)

func labmodel(c color.Color) color.Color {
	if _, ok := c.(Color); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	l, la, lb := FromRGB(r, g, b)
	return Color{l, la, lb}
}

type Color struct {
	L, A, B float64
}

func (c Color) RGBA() (r, g, b, a uint32) {
	r, g, b = ToRGB(c.L, c.A, c.B)
	return r, g, b, 0xffff
}

type Image struct {
	// L, A, B hold the appropriate image coordinates
	L, A, B []float64
	// Stride is the L, A and B stride in pixels
	Stride int
	// Image bounds
	Rect image.Rectangle
}

func NewImage(r image.Rectangle) *Image {
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

func (p *Image) ColorModel() color.Model { return color.NRGBA64Model }

func (p *Image) Bounds() image.Rectangle { return p.Rect }

func (p *Image) At(x, y int) color.Color { return p.LabAt(x, y) }

func (p *Image) LabAt(x, y int) Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return Color{}
	}
	i := p.PixOffset(x, y)
	return Color{p.L[i], p.A[i], p.B[i]}
}

func (p *Image) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x - p.Rect.Min.X)
}

func (p *Image) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}

	i := p.PixOffset(x, y)
	c1 := Model.Convert(c).(Color)
	p.L[i], p.A[i], p.B[i] = c1.L, c1.A, c1.B
}

func (p *Image) SetLab(x, y int, c Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}

	i := p.PixOffset(x, y)
	p.L[i], p.A[i], p.B[i] = c.L, c.A, c.B
}
