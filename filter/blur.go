package filter

import "github.com/loov/sketch-capture/lab"

func avg(a, b, c float64) float64 {
	return (a + b + c) / 3
}

func BoxBlur(m *lab.Image, steps int) {
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
