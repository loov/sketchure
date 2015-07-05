package filter

import "github.com/loov/sketch-capture/cielab"

func Desaturate(m *cielab.Image) {
	m.A = make([]float64, len(m.A))
	m.B = make([]float64, len(m.B))
}