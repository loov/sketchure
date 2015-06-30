package filter

import (
	"image"
	"math/rand"

	"github.com/loov/sketch-capture/lab"
)

func ValueNoise(m *lab.Image, strength float64) {
	max := 100 * strength
	for i := range m.L {
		m.L[i] += (rand.Float64() - 0.5) * max
		if m.L[i] < 0 {
			m.L[i] = 0
		} else if m.L[i] > 100 {
			m.L[i] = 100
		}

	}
}

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
func NormalizeGradientFromCorners(m *lab.Image) {
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

func NormalizeGradient(m *lab.Image) {
	Desaturate(m)

	r := m.Bounds()

	base := m.Clone()
	Erode(base, 5)
	Blur(base, 5)

	white := 100.0
	average := average(m, base.Bounds())
	invspan := 1 / (average / 100)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := m.Offset(r.Min.X, y)
		for x := r.Min.X; x < r.Max.X; x++ {
			v := m.L[i] - base.L[i] + white
			x := 100 - (100-v)*invspan
			m.L[i] = x
			i++
		}
	}
}

func max(a, b, c float64) float64 {
	if a > b {
		if a > c {
			return a
		} else {
			return c
		}
	} else {
		if b > c {
			return b
		} else {
			return c
		}
	}
}

func Erode(m *lab.Image, steps int) {
	for i := 0; i < steps; i++ {
		ErodeHorizontal3(m)
		ErodeVertical3(m)
	}
}

func ErodeHorizontal3(m *lab.Image) {
	r := m.Bounds()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := m.Offset(r.Min.X, y)
		p := m.L[i]
		for x := r.Min.X; x < r.Max.X-1; x++ {
			z, n := m.L[i], m.L[i+1]
			m.L[i] = max(p, z, n)
			p = z
			i++
		}
		m.L[i] = max(p, m.L[i], m.L[i])
	}
}

func ErodeVertical3(m *lab.Image) {
	r := m.Bounds()
	stride := m.Offset(r.Min.X, r.Min.Y+1) - m.Offset(r.Min.X, r.Min.Y)
	for x := r.Min.X; x < r.Max.X; x++ {
		i := m.Offset(x, r.Min.Y)
		p := m.L[i]
		for y := r.Min.Y; y < r.Max.Y-1; y++ {
			z, n := m.L[i], m.L[i+stride]
			m.L[i] = max(p, z, n)
			p = z
			i += stride
		}
		m.L[i] = max(p, m.L[i], m.L[i])
	}
}

func avg(a, b, c float64) float64 {
	return (a + b + c) / 3
}

func Blur(m *lab.Image, steps int) {
	for i := 0; i < steps; i++ {
		BlurHorizontal3(m)
		BlurVertical3(m)
	}
}

func BlurHorizontal3(m *lab.Image) {
	r := m.Bounds()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := m.Offset(r.Min.X, y)
		p := m.L[i]
		for x := r.Min.X; x < r.Max.X-1; x++ {
			z, n := m.L[i], m.L[i+1]
			m.L[i] = avg(p, z, n)
			p = z
			i++
		}
		m.L[i] = avg(p, m.L[i], m.L[i])
	}
}

func BlurVertical3(m *lab.Image) {
	r := m.Bounds()
	stride := m.Offset(r.Min.X, r.Min.Y+1) - m.Offset(r.Min.X, r.Min.Y)
	for x := r.Min.X; x < r.Max.X; x++ {
		i := m.Offset(x, r.Min.Y)
		p := m.L[i]
		for y := r.Min.Y; y < r.Max.Y-1; y++ {
			z, n := m.L[i], m.L[i+stride]
			m.L[i] = avg(p, z, n)
			p = z
			i += stride
		}
		m.L[i] = avg(p, m.L[i], m.L[i])
	}
}
