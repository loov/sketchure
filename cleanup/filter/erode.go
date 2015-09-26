package filter

func max(a, b, c float32) float32 {
	if a >= b {
		if a >= c {
			return a
		}
	} else if b > c {
		return b
	}
	return c
}

// Erode channel with 3x3 kernel
func (ch *Channel) Erode(steps int) {
	for i := 0; i < steps; i++ {
		ch.ErodeH3()
		ch.ErodeV3()
	}
}

// Erode channel horizontally with 3px kernel
func (ch *Channel) ErodeH3() {
	data, w, h, stride := ch.Data, ch.Width, ch.Height, ch.Stride

	for y := 0; y < h; y++ {
		i := y * stride
		e := y*stride + w - 1
		p, z := data[i], data[i]
		for ; i < e; i++ {
			n := data[i+1]
			data[i] = max(p, z, n)
			p, z = z, n
		}
		data[i] = max(p, data[i], data[i])
	}
}

// Erode channel vertically with 3px kernel
func (ch *Channel) ErodeV3() {
	data, w, h, stride := ch.Data, ch.Width, ch.Height, ch.Stride

	for x := 0; x < w; x++ {
		i := x
		e := (h-1)*stride + x
		p, z := data[i], data[i]
		for ; i < e; i += stride {
			n := data[i+stride]
			data[i] = max(p, z, n)
			p, z = z, n
		}
		data[i] = max(p, data[i], data[i])
	}
}
