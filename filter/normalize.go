package filter

import (
	"image"

	"github.com/loov/sketch-capture/lab"
)

func Desaturate(m *lab.Image) {
	m.A = make([]float64, len(m.A))
	m.B = make([]float64, len(m.B))
}

func lerp(c, min, max int, minval, maxval float64) float64 {
	p := float64(c-min) / float64(max-min)
	return p*(maxval-minval) + minval
}

//TODO: optimize
func average(m *lab.Image, r image.Rectangle) (L float64) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := m.Offset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			L += m.L[i]
			i++
		}
	}

	return L / float64(r.Dx()*r.Dy())
}

//TODO: optimize
func NormalizeCornersBilinear(m *lab.Image) {
	const e = 20 // corner sample size

	r := m.Bounds()

	topLeft := average(m, image.Rect(r.Min.X, r.Min.Y, r.Min.X+e, r.Min.Y+e))
	topRight := average(m, image.Rect(r.Max.X-e, r.Min.Y, r.Max.X, r.Min.Y+e))
	bottomLeft := average(m, image.Rect(r.Min.X, r.Max.Y-e, r.Min.X+e, r.Max.Y))
	bottomRight := average(m, image.Rect(r.Max.X-e, r.Max.Y-e, r.Max.X, r.Max.Y))

	white := 100.0
	average := (topLeft + topRight + bottomLeft + bottomRight) / 4
	invspan := 1 / (average / 100)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		left := lerp(y, r.Min.Y, r.Max.Y, topLeft, bottomLeft)
		right := lerp(y, r.Min.Y, r.Max.Y, topRight, bottomRight)

		grad := (right - left) / float64(r.Max.X-r.Min.X)
		base := left

		i := m.Offset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			v := m.L[i] - base + white
			x := 100 - (100-v)*invspan
			m.L[i] = x
			base += grad
			i++
		}
	}
}

func Normalize(m *lab.Image) {
	const white = 110.0

	Desaturate(m)

	r := m.Bounds()

	base := m.Clone()
	Erode(base, 15)
	BoxBlur(base, 15)

	average := average(m, base.Bounds())
	invspan := 1 / (average / white)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := m.Offset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			v := m.L[i] - base.L[i] + white
			x := white - (white-v)*invspan
			m.L[i] = x
			i++
		}
	}
}
