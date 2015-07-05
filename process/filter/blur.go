package filter

import "github.com/loov/sketch-capture/cielab"

func avg(a, b, c float64) float64 {
	return (a + b + c) / 3
}

// Blur L channel with 3x3 kernel
func Blur(m *cielab.Image, steps int) {
	for i := 0; i < steps; i++ {
		BlurHorizontal3(m)
		BlurVertical3(m)
	}
}

// Blur L channel horizontally with 3px kernel
func BlurHorizontal3(m *cielab.Image) {
	r := m.Bounds()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := m.Offset(r.Min.X, y)
		p := m.L[i]
		z := p
		for x := r.Min.X; x < r.Max.X-1; x++ {
			n := m.L[i+1]
			m.L[i] = avg(p, z, n)
			p, z = z, n
			i++
		}
		m.L[i] = avg(p, m.L[i], m.L[i])
	}
}

// Blur L channel vertically with 3px kernel
func BlurVertical3(m *cielab.Image) {
	r := m.Bounds()
	stride := m.Offset(r.Min.X, r.Min.Y+1) - m.Offset(r.Min.X, r.Min.Y)
	for x := r.Min.X; x < r.Max.X; x++ {
		i := m.Offset(x, r.Min.Y)
		p := m.L[i]
		z := p
		for y := r.Min.Y; y < r.Max.Y-1; y++ {
			n := m.L[i+stride]
			m.L[i] = avg(p, z, n)
			p, z = z, n
			i += stride
		}
		m.L[i] = avg(p, m.L[i], m.L[i])
	}
}
