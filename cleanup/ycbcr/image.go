package ycbcr

import (
	"image"
	"image/color"
)

// Image represents an image in YCbCr colorspace with float32 precision
type Image struct {
	// Y  is in range [0..1]
	// Cb is in range [-1..1]
	// Cr is in range [-1..1]
	Y, Cb, Cr []float32
	Stride    int
	Rect      image.Rectangle
}



func NewImage(r image.Rectangle) *Image {
	w, h := r.Dx(), r.Dy()
	sz := w * h * 4
	return &Image{
		Y:  make([]float32, sz),
		Cb: make([]float32, sz),
		Cr: make([]float32, sz),

		Stride: w,
		Rect:   r,
	}
}

func (m *Image) Clone() *Image {
	cp := &Image{}
	*cp = *m

	cp.Y = make([]float32, len(m.Y))
	copy(cp.Y, m.Y)
	cp.Cb = make([]float32, len(m.Cb))
	copy(cp.Cb, m.Cb)
	cp.Cr = make([]float32, len(m.Cr))
	copy(cp.Cr, m.Cr)

	return cp
}

func (m *Image) ColorModel() color.Model { return ColorModel }

func (m *Image) Bounds() image.Rectangle { return m.Rect }

func (m *Image) At(x, y int) color.Color { return m.YCbCrAt(x, y) }

func (m *Image) YCbCrAt(x, y int) Color {
	if !(image.Point{x, y}.In(m.Rect)) {
		return Color{}
	}
	i := m.Offset(x, y)
	return Color{m.Y[i], m.Cb[i], m.Cr[i]}
}

func (m *Image) Offset(x, y int) int {
	return (y-m.Rect.Min.Y)*m.Stride + (x - m.Rect.Min.X)
}

func (m *Image) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(m.Rect)) {
		return
	}

	i := m.Offset(x, y)
	c1 := ColorModel.Convert(c).(Color)
	m.Y[i], m.Cb[i], m.Cr[i] = c1.Y, c1.Cb, c1.Cr
}

func (m *Image) SetLAB(x, y int, c Color) {
	if !(image.Point{x, y}.In(m.Rect)) {
		return
	}

	i := m.Offset(x, y)
	m.Y[i], m.Cb[i], m.Cr[i] = c.Y, c.Cb, c.Cr
}