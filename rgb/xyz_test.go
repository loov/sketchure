package rgb

import (
	"math/rand"
	"testing"
)

func closeto(got, exp uint32) bool {
	var d uint32
	if got > exp {
		d = got - exp
	} else {
		d = exp - got
	}
	return d < 0x20
}

func randRGB() (uint32, uint32, uint32) {
	return rand.Uint32() & 0xFFFF, rand.Uint32() & 0xFFFF, rand.Uint32() & 0xFFFF
}

func TestConvertXYZ(t *testing.T) {
	var x, y, z float64
	var r, g, b uint32

	x, y, z = ToXYZ(0, 0, 0)
	r, g, b = FromXYZ(x, y, z)

	if !closeto(r, 0) || !closeto(g, 0) || !closeto(b, 0) {
		t.Errorf("expected <0000,0000,0000> got <%4x,%4x,%4x>", r, g, b)
	}

	x, y, z = ToXYZ(0xFFFF, 0xFFFF, 0xFFFF)
	r, g, b = FromXYZ(x, y, z)
	if !closeto(r, 0xFFFF) || !closeto(g, 0xFFFF) || !closeto(b, 0xFFFF) {
		t.Errorf("expected <ffff,ffff,ffff> got <%4x,%4x,%4x>", r, g, b)
	}

	for i := 0; i < 100; i++ {
		er, eg, eb := randRGB()
		x, y, z = ToXYZ(er, eg, eb)
		r, g, b = FromXYZ(x, y, z)
		if !closeto(r, er) || !closeto(g, eg) || !closeto(b, eb) {
			t.Errorf("expected <%4x,%4x,%4x> got <%4x,%4x,%4x>", er, eg, eb, r, g, b)
		}
	}
}
