package cielab

import "testing"

func TestConvertXYZToLAB(t *testing.T) {
	const eps = 1e-14

	for i := 0; i < 1<<10; i++ {
		ex, ey, ez := rand3()
		l, a, b := XYZToLAB(ex, ey, ez)
		x, y, z := LABToXYZ(l, a, b)
		if !closeto(x, ex, eps) || !closeto(y, ey, eps) || !closeto(z, ez, eps) {
			t.Fatalf("expected <%.4f,%.4f,%.4f> got <%.4f,%.4f,%.4f>", ex, ey, ez, x, y, z)
		}
	}
}
