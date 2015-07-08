package cielab

import "github.com/loov/sketchure/cielab/gamma"

func SRGB16ToLAB(r, g, b uint32) (l, la, lb float64) {
	rl, gl, bl := gamma.LinearizeRGB16(uint16(r), uint16(g), uint16(b))
	x, y, z := RGBToXYZ(rl, gl, bl)
	l, la, lb = XYZToLAB(x, y, z)
	return
}

func LABToSRGB16(l, la, lb float64) (r, g, b uint32) {
	x, y, z := LABToXYZ(l, la, lb)
	rf, gf, bf := XYZToRGB(x, y, z)
	r2, g2, b2 := gamma.DelinearizeRGB16(rf, gf, bf)
	return uint32(r2), uint32(g2), uint32(b2)
}
