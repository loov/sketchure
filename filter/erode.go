package filter

import "github.com/loov/sketch-capture/cielab"

func max(a, b, c float64) float64 {
	if a >= b {
		if a >= c {
			return a
		}
	} else if b > c {
		return b
	}
	return c
}

func Erode(m *cielab.Image, steps int) {
	for i := 0; i < steps; i++ {
		ErodeHorizontal3(m)
		ErodeVertical3(m)
	}
}

func ErodeHorizontal3(m *cielab.Image) {
	r := m.Bounds()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := m.Offset(r.Min.X, y)
		p := m.L[i]
		z := p
		for x := r.Min.X; x < r.Max.X-1; x++ {
			n := m.L[i+1]
			m.L[i] = max(p, z, n)
			p, z = z, n
			i++
		}
		m.L[i] = max(p, m.L[i], m.L[i])
	}
}

func ErodeVertical3(m *cielab.Image) {
	r := m.Bounds()
	stride := m.Offset(r.Min.X, r.Min.Y+1) - m.Offset(r.Min.X, r.Min.Y)
	for x := r.Min.X; x < r.Max.X; x++ {
		i := m.Offset(x, r.Min.Y)
		p := m.L[i]
		z := p
		for y := r.Min.Y; y < r.Max.Y-1; y++ {
			n := m.L[i+stride]
			m.L[i] = max(p, z, n)
			p, z = z, n
			i += stride
		}
		m.L[i] = max(p, m.L[i], m.L[i])
	}
}
