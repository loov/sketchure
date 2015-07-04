package cielab

import "math/rand"

func u16tof(v uint32) float64 {
	return float64(v) / 0xFFFF
}

func ftou16(v float64) uint32 {
	if v < 0 {
		return 0
	} else if v > 1.0 {
		return 0xFFFF
	}
	return uint32(v * 0xFFFF)
}

func u16_f64(a, b, c uint32) (af, bf, cf float64) {
	return u16tof(a), u16tof(b), u16tof(c)
}

func f64_u16(af, bf, cf float64) (a, b, c uint32) {
	return ftou16(af), ftou16(bf), ftou16(cf)
}

func closeto(got, exp, eps float64) bool {
	var d float64
	if got > exp {
		d = got - exp
	} else {
		d = exp - got
	}
	return d < eps
}

func rand3() (float64, float64, float64) {
	return rand.Float64(), rand.Float64(), rand.Float64()
}
