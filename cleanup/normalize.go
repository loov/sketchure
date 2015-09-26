package cleanup

import (
	"github.com/loov/sketchure/cleanup/filter"
	"github.com/loov/sketchure/cleanup/ycbcr"
)

func lerp(c, min, max int, minval, maxval float64) float64 {
	p := float64(c-min) / float64(max-min)
	return p*(maxval-minval) + minval
}

type Options struct {
	// Whiteness determines the lightest value in output image
	//   normal range is 0..1
	// if you wish to cut off light values use value larger than 1
	Whiteness float32

	// LineWidth is the maximum line width on the page in pixels
	// The higher this value is, the less tolerant the algorithm
	//   is to quick gradient/lighting changes
	// If the value is too small, there will be errors in solid areas
	LineWidth int
}

func DefaultOptionsFor(m *ycbcr.Image) *Options {
	r := m.Bounds()
	return &Options{
		Whiteness: 100,
		LineWidth: r.Dx() / 20,
	}
}

// ByBase cleans image based on the surrounding color values
func ByBase(m *ycbcr.Image, opts *Options) {
	if opts == nil {
		opts = DefaultOptionsFor(m)
	}

	white := opts.Whiteness

	L := &filter.Channel{
		Data:   m.Y,
		Width:  m.Bounds().Dx(),
		Height: m.Bounds().Dy(),
		Stride: m.Stride,
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
			L.Data[i] = white + (L.Data[i]-base.Data[i])*invspan
		}
	}
}
