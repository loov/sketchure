package ycbcr

import (
	"image/color"

	"github.com/loov/sketchure/cleanup/gamma"
)

// Color represents a color in LAB
type Color struct{ Y, Cb, Cr float32 }

var ColorModel = color.ModelFunc(convert)

func convert(c color.Color) color.Color {
	if _, ok := c.(Color); ok {
		return c
	}

	r, g, b, _ := c.RGBA()
	yy, cb, cr := SRGB16ToYCbCr(r, g, b)
	return Color{yy, cb, cr}
}

func (c Color) RGBA() (r, g, b, a uint32) {
	r, g, b = YCbCrToSRGB16(c.Y, c.Cb, c.Cr)
	return r, g, b, 0xffff
}

func SRGB16ToYCbCr(r, g, b uint32) (Y, Cb, Cr float32) {
	R, G, B := gamma.LinearizeRGB16(uint16(r), uint16(g), uint16(b))

	Y = 0.2990*R + 0.5870*G + 0.1140*B
	Cb = -0.1687*R - 0.3313*G + 0.5000*B
	Cr = 0.5000*R - 0.4187*G - 0.0813*B

	return
}

func YCbCrToSRGB16(Y, Cb, Cr float32) (r, g, b uint32) {
	R := Y + 1.40200*Cr
	G := Y - 0.34414*Cb - 0.71414*Cr
	B := Y + 1.77200*Cb

	r2, g2, b2 := gamma.DelinearizeRGB16(R, G, B)
	return uint32(r2), uint32(g2), uint32(b2)
}
