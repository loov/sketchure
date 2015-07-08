package cleanup

import (
	"image"

	"github.com/loov/sketchure/cielab"
	"github.com/loov/sketchure/cleanup/filter"
)

func lerp(c, min, max int, minval, maxval float64) float64 {
	p := float64(c-min) / float64(max-min)
	return p*(maxval-minval) + minval
}

//TODO: optimize
func average(m *cielab.Image, r image.Rectangle) (L float64) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := m.Offset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			L += m.L[i]
			i++
		}
	}
	return L / float64(r.Dx()*r.Dy())
}

type Options struct {
	// Corner size in pixels to sample for determining base level
	CornerSize int

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

func DefaultOptionsFor(m *cielab.Image) *Options {
	r := m.Bounds()
	return &Options{
		Whiteness:  105,
		CornerSize: r.Dx() / 20,
		LineWidth:  r.Dx() / 20,
	}
}

// ByCorners cleans image based on the image corners
func ByCorners(m *cielab.Image, opts *Options) {
	r := m.Bounds()

	if opts == nil {
		opts = DefaultOptionsFor(m)
	}
	cs := opts.CornerSize

	topLeft := average(m, image.Rect(r.Min.X, r.Min.Y, r.Min.X+cs, r.Min.Y+cs))
	topRight := average(m, image.Rect(r.Max.X-cs, r.Min.Y, r.Max.X, r.Min.Y+cs))
	bottomLeft := average(m, image.Rect(r.Min.X, r.Max.Y-cs, r.Min.X+cs, r.Max.Y))
	bottomRight := average(m, image.Rect(r.Max.X-cs, r.Max.Y-cs, r.Max.X, r.Max.Y))

	white := opts.Whiteness
	average := (topLeft + topRight + bottomLeft + bottomRight) / 4
	invspan := 1 / (average / 100)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		left := lerp(y, r.Min.Y, r.Max.Y, topLeft, bottomLeft)
		right := lerp(y, r.Min.Y, r.Max.Y, topRight, bottomRight)

		grad := (right - left) / float64(r.Max.X-r.Min.X)
		base := left

		i := m.Offset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			m.L[i] = white + (m.L[i]-base)*invspan
			base += grad
			i++
		}
	}
}

// ByBase cleans image based on the surrounding color values
func ByBase(m *cielab.Image, opts *Options) {
	r := m.Bounds()

	if opts == nil {
		opts = DefaultOptionsFor(m)
	}

	white := opts.Whiteness

	base := m.Clone()
	filter.Erode(base, opts.LineWidth)
	filter.Blur(base, opts.LineWidth)

	average := average(m, base.Bounds())
	invspan := 1 / (average / white)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := m.Offset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			m.L[i] = white + (m.L[i]-base.L[i])*invspan
			i++
		}
	}
}
