package filter

func avg(a, b, c float32) float32 {
	return (a + b + c) / 3
}

// Blur channel with 3x3 kernel
func (ch *Channel) Blur(steps int) {
	for i := 0; i < steps; i++ {
		ch.BlurH3()
		ch.BlurV3()
	}
}

// Blur channel horizontally with 3px kernel
func (ch *Channel) BlurH3() {
	data, w, h, stride := ch.Data, ch.Width, ch.Height, ch.Stride

	for y := 0; y < h; y++ {
		i := y * stride
		e := y*stride + w - 1
		p, z := data[i], data[i]
		for ; i < e; i++ {
			n := data[i+1]
			data[i] = avg(p, z, n)
			p, z = z, n
		}
		data[i] = avg(p, data[i], data[i])
	}
}

// Blur channel vertically with 3px kernel
func (ch *Channel) BlurV3() {
	data, w, h, stride := ch.Data, ch.Width, ch.Height, ch.Stride

	for x := 0; x < w; x++ {
		i := x
		e := (h-1)*stride + x
		p, z := data[i], data[i]
		for ; i < e; i += stride {
			n := data[i+stride]
			data[i] = avg(p, z, n)
			p, z = z, n
		}
		data[i] = avg(p, data[i], data[i])
	}
}
