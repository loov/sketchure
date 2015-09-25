package filter

import "image"

func Desaturate(m *image.YCbCr) {
	for i := 0; i < len(m.Cb); i++ {
		m.Cb[i] = 127
	}
	for i := 0; i < len(m.Cr); i++ {
		m.Cr[i] = 127
	}
}

type Channel struct {
	Data   []byte
	Width  int
	Height int
	Stride int
}

func (ch *Channel) Average() (avg float64) {
	data, w, h, stride := ch.Data, ch.Width, ch.Height, ch.Stride

	for y := 0; y < h; y++ {
		i := y * stride
		e := i + w
		for ; i < e; i++ {
			avg += float64(data[i])
		}
	}

	return avg / float64(w*h)
}

func (ch *Channel) Clone() *Channel {
	cp := *ch
	cp.Data = make([]byte, len(ch.Data))
	copy(cp.Data, ch.Data)
	return &cp
}
