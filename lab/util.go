package lab

import (
	"image"
	"image/draw"
	"os"
)

func ImageFrom(src image.Image) *Image {
	r := src.Bounds()
	if r.Min.X == 0 && r.Min.Y == 0 {
		if n, ok := src.(*Image); ok {
			return n
		}
	}

	dst := NewImage(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	return dst
}

func ImageFromFile(filename string) (*Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return ImageFrom(m), nil
}
