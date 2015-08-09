package filter

import "github.com/loov/sketchure/cielab"

func mid(a, b, c float64) float64 {
	if a > b {
		if b > c {
			return b
		} else if a < c {
			return a
		}
	} else {
		if a > c {
			return a
		} else if b < c {
			return b
		}
	}
	return c
}

// Median L channel with 3x3 kernel
func Median(m *cielab.Image, steps int) {
	for i := 0; i < steps; i++ {
		MedianHorizontal3(m)
		MedianVertical3(m)
	}
}

// Median L channel horizontally with 3px kernel
func MedianHorizontal3(m *cielab.Image) {
	r := m.Bounds()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		i := m.Offset(r.Min.X, y)
		p := m.L[i]
		z := p
		for x := r.Min.X; x < r.Max.X-1; x++ {
			n := m.L[i+1]
			m.L[i] = mid(p, z, n)
			p, z = z, n
			i++
		}
		m.L[i] = mid(p, m.L[i], m.L[i])
	}
}

// Median L channel vertically with 3px kernel
func MedianVertical3(m *cielab.Image) {
	r := m.Bounds()
	stride := m.Offset(r.Min.X, r.Min.Y+1) - m.Offset(r.Min.X, r.Min.Y)
	for x := r.Min.X; x < r.Max.X; x++ {
		i := m.Offset(x, r.Min.Y)
		p := m.L[i]
		z := p
		for y := r.Min.Y; y < r.Max.Y-1; y++ {
			n := m.L[i+stride]
			m.L[i] = mid(p, z, n)
			p, z = z, n
			i += stride
		}
		m.L[i] = mid(p, m.L[i], m.L[i])
	}
}
