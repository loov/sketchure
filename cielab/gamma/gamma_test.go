package gamma

import (
	"math"
	"math/rand"
	"testing"
)

func TestValueConversion(t *testing.T) {
	const epsilon = 1e-14
	errors := 0
	for i := 0; i < 1<<16; i++ {
		v := rand.Float64()
		v2 := Linearize(v)
		vx := Delinearize(v2)
		d := math.Abs(v - vx)
		if d > epsilon {
			t.Errorf("delta %f : %f => %f => %f", d, v, v2, vx)
			errors++
			if errors > 10 {
				break
			}
		}
	}
}
