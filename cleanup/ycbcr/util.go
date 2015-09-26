package ycbcr

import (
	"image"
	"os"
)

func FromFile(filename string) (*Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return FromImage(src), nil
}

func FromImage(src image.Image) *Image {
	r := src.Bounds()
	dst := NewImage(src.Bounds())
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			c := src.At(x, y)
			r, g, b, _ := c.RGBA()
			yy, cb, cr := SRGB16ToYCbCr(r, g, b)

			i := dst.Offset(x, y)
			dst.Y[i] = yy
			dst.Cb[i] = cb
			dst.Cr[i] = cr
		}
	}
	return dst
}
