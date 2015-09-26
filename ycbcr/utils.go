package ycbcr

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

func FromFile(filename string) (*image.YCbCr, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	if ycbcr, ok := src.(*image.YCbCr); ok {
		return ycbcr, nil
	}

	return FromImage(src), nil
}

func ToFile(filename string, m image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if filepath.Ext(filename) == ".png" {
		return png.Encode(file, m)
	}
	return jpeg.Encode(file, m, &jpeg.Options{Quality: 80})
}

func FromImage(src image.Image) *image.YCbCr {
	r := src.Bounds()
	dst := image.NewYCbCr(src.Bounds(), image.YCbCrSubsampleRatio444)
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			c := src.At(x, y)
			r, g, b, _ := c.RGBA()
			yy, cb, cr := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))

			i := dst.YOffset(x, y)
			dst.Y[i] = yy
			dst.Cb[i] = cb
			dst.Cr[i] = cr
		}
	}

	return dst
}

func Clone(src *image.YCbCr) *image.YCbCr {
	dst := *src
	dst.Y = make([]uint8, len(src.Y))
	copy(dst.Y, src.Y)
	dst.Cb = make([]uint8, len(src.Cb))
	copy(dst.Cb, src.Cb)
	dst.Cr = make([]uint8, len(src.Cr))
	copy(dst.Cr, src.Cr)
	return &dst
}
