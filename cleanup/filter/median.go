package filter

func mid(a, b, c byte) byte {
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

// Median channel with 3x3 kernel
func (ch *Channel) Median(steps int) {
	for i := 0; i < steps; i++ {
		ch.MedianH3()
		ch.MedianV3()
	}
}

// Median channel horizontally with 3px kernel
func (ch *Channel) MedianH3() {
	data, w, h, stride := ch.Data, ch.Width, ch.Height, ch.Stride

	for y := 0; y < h; y++ {
		i := y * stride
		e := y*stride + w - 1
		p, z := data[i], data[i]
		for ; i < e; i++ {
			n := data[i+1]
			data[i] = mid(p, z, n)
			p, z = z, n
		}
		data[i] = mid(p, data[i], data[i])
	}
}

// Median channel vertically with 3px kernel
func (ch *Channel) MedianV3() {
	data, w, h, stride := ch.Data, ch.Width, ch.Height, ch.Stride

	for x := 0; x < w; x++ {
		i := x
		e := (h-1)*stride + x
		p, z := data[i], data[i]
		for ; i < e; i += stride {
			n := data[i+stride]
			data[i] = mid(p, z, n)
			p, z = z, n
		}
		data[i] = mid(p, data[i], data[i])
	}
}
