package cielab

import (
	"math/rand"
	"testing"
)

func r16() uint32 {
	return rand.Uint32() & 0xFFFF
}

func TestConvertRGBToXYZ(t *testing.T) {
	const eps = 0xC

	for i := 0; i < 1<<10; i++ {
		er, eg, eb := r16(), r16(), r16()
		l, la, lb := SRGB16ToLAB(er, eg, eb)
		r, g, b := LABToSRGB16(l, la, lb)
		if !closeto16(r, er, eps) || !closeto16(g, eg, eps) || !closeto16(b, eb, eps) {
			t.Fatalf("expected <%.4f,%.4f,%.4f> got <%.4f,%.4f,%.4f>", er, eg, eb, r, g, b)
		}
	}
}

func closeto16(got, exp, eps uint32) bool {
	var d uint32
	if got > exp {
		d = got - exp
	} else {
		d = exp - got
	}
	return d < eps
}
