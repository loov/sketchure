package gamma

import "math"

// Linearizes value based on sRGB gamma
func Linearize(v float64) float64 {
	if v < 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// Delinearizes value based on sRGB gamma
func Delinearize(v float64) float64 {
	if v < 0.0031308 {
		return v * 12.92
	}
	return 1.055*math.Pow(v, 1.0/2.4) - 0.055
}

// Converts from sRGB to linear RGB
func LinearizeRGB(r, g, b float64) (r2, g2, b2 float64) {
	return Linearize(r), Linearize(g), Linearize(b)
}

// Converts from linear RGB to sRGB
func DelinearizeRGB(r2, g2, b2 float64) (r, g, b float64) {
	return Delinearize(r2), Delinearize(g2), Delinearize(b2)
}

var linearizeRGB16 = [1 << 16]float64{}

func init() {
	const max = 0xFFFF
	for i := range linearizeRGB16 {
		v := float64(i) / float64(max)
		linearizeRGB16[i] = Linearize(v)
	}
}

// Converts from sRGB to linear RGB
func LinearizeRGB16(r, g, b uint16) (r2, g2, b2 float64) {
	r2 = linearizeRGB16[r]
	g2 = linearizeRGB16[g]
	b2 = linearizeRGB16[b]
	return
}

// Converts from linear RGB to sRGB
func DelinearizeRGB16(r2, g2, b2 float64) (r, g, b uint16) {
	r = u16(Delinearize(r2))
	g = u16(Delinearize(g2))
	b = u16(Delinearize(b2))
	return
}

func u16(v float64) uint16 {
	if v <= 0 {
		return 0
	} else if v >= 1 {
		return 0xFFFF
	}
	return uint16(v * 0xFFFF)
}
