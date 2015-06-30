package rgb

import (
	"image"
	"image/draw"
)

func ImageFrom(src image.Image) *image.NRGBA {
	r := src.Bounds()
	if r.Min.X == 0 && r.Min.Y == 0 {
		if n, ok := src.(*image.NRGBA); ok {
			return n
		}
	}

	dst := image.NewNRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	return dst
}
