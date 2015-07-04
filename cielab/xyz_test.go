package cielab

import "testing"

func TestConvertRGBToLAB(t *testing.T) {
	const eps = 1e-6

	for i := 0; i < 1<<10; i++ {
		er, eg, eb := rand3()
		x, y, z := RGBToXYZ(er, eg, eb)
		r, g, b := XYZToRGB(x, y, z)
		if !closeto(r, er, eps) || !closeto(g, eg, eps) || !closeto(b, eb, eps) {
			t.Fatalf("expected <%.4f,%.4f,%.4f> got <%.4f,%.4f,%.4f>", er, eg, eb, r, g, b)
		}
	}
}
