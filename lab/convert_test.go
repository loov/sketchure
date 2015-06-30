package lab

import (
	"math/rand"
	"testing"
)

func closeto(got, exp float64) bool {
	var d float64
	if got > exp {
		d = got - exp
	} else {
		d = exp - got
	}
	return d < 0.1/255.0
}

func randXYZ() (x, y, z float64) {
	return rand.Float64(), rand.Float64(), rand.Float64()
}

func TestConvertXYZ(t *testing.T) {
	for i := 0; i < 100; i++ {
		ex, ey, ez := randXYZ()
		l, a, b := FromXYZ(ex, ey, ez)
		x, y, z := ToXYZ(l, a, b)
		if !closeto(x, ex) || !closeto(y, ey) || !closeto(z, ez) {
			t.Fatalf("expected <%.4f,%.4f,%.4f> got <%.4f,%.4f,%.4f>", ex, ey, ez, x, y, z)
		}
	}
}
