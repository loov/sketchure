package cleanup

import (
	"image"

	"github.com/loov/sketchure/cleanup/filter"
)

func lerp(c, min, max int, minval, maxval float64) float64 {
	p := float64(c-min) / float64(max-min)
	return p*(maxval-minval) + minval
}

type Options struct {
	// Whiteness determines the lightest value in output image
	//   normal range is 0..100
	// if you wish to cut off light values use value larger than 100
	Whiteness float64

	// LineWidth is the maximum line width on the page in pixels
	// The higher this value is, the less tolerant the algorithm
	//   is to quick gradient/lighting changes
	// If the value is too small, there will be errors in solid areas
	LineWidth int
}

func DefaultOptionsFor(m *image.YCbCr) *Options {
	r := m.Bounds()
	return &Options{
		Whiteness: 100,
		LineWidth: r.Dx() / 20,
	}
}

// ByBase cleans image based on the surrounding color values
func ByBase(m *image.YCbCr, opts *Options) {
	if opts == nil {
		opts = DefaultOptionsFor(m)
	}

	white := opts.Whiteness * 255.0 / 100.0

	L := &filter.Channel{
		Data:   m.Y,
		Width:  m.Bounds().Dx(),
		Height: m.Bounds().Dy(),
		Stride: m.YStride,
	}

	// get rid of hot-pixels
	L.Median(1)

	base := L.Clone()
	base.Erode(opts.LineWidth)
	base.Blur(opts.LineWidth)

	average := base.Average()
	invspan := 1.0 / (average / white)

	for y := 0; y < L.Height; y++ {
		i := y * L.Stride
		e := i + L.Width
		for ; i < e; i++ {
			lv := float64(L.Data[i])
			bv := float64(base.Data[i])

			r := int(white + (lv-bv)*invspan)
			if r < 0x00 {
				r = 0x00
			} else if r > 0xFF {
				r = 0xFF
			}

			L.Data[i] = byte(r)
		}
	}
}
