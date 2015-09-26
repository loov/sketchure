package filter

import "github.com/loov/sketchure/cleanup/ycbcr"

func Desaturate(m *ycbcr.Image) {
	for i := 0; i < len(m.Cb); i++ {
		m.Cb[i] = 0.0
	}
	for i := 0; i < len(m.Cr); i++ {
		m.Cr[i] = 0.0
	}
}

type Channel struct {
	Data   []float32
	Width  int
	Height int
	Stride int
}

func (ch *Channel) Average() (avg float32) {
	data, w, h, stride := ch.Data, ch.Width, ch.Height, ch.Stride

	for y := 0; y < h; y++ {
		i := y * stride
		e := i + w
		for ; i < e; i++ {
			avg += data[i]
		}
	}

	return avg / float32(w*h)
}

func (ch *Channel) Clone() *Channel {
	cp := *ch
	cp.Data = make([]float32, len(ch.Data))
	copy(cp.Data, ch.Data)
	return &cp
}
