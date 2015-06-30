package rgb

import "math"

func lF64(v float64) float64 {
	if v < 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func dF64(v float64) float64 {
	if v < 0.0031308 {
		return v * 12.92
	}
	return 1.055*math.Pow(v, 1.0/2.4) - 0.055
}

func LinearizeF64(r, g, b float64) (float64, float64, float64) {
	return lF64(r), lF64(g), lF64(b)
}

func DelinearizeF64(r, g, b float64) (float64, float64, float64) {
	return dF64(r), dF64(g), dF64(b)
}

var (
	lU32 = [1 << 16]uint32{}
	dU32 = [1 << 16]uint32{}
)

func init() {
	for i := range lU32 {
		lU32[i] = f2u(lF64(u2f(uint32(i))))
		dU32[i] = f2u(dF64(u2f(uint32(i))))
	}
}

// r, g, b must be in range [0, 0xFFFF]
func LinearizeU32(r, g, b uint32) (uint32, uint32, uint32) {
	return lU32[r], lU32[g], lU32[b]
}

// r, g, b must be in range [0, 0xFFFF]
func DelinearizeU32(r, g, b uint32) (uint32, uint32, uint32) {
	return dU32[r], dU32[g], dU32[b]
}

func u2f(v uint32) float64 { return float64(v) / 0xFFFF }
func f2u(v float64) uint32 {
	if v < 0 {
		return 0
	} else if v > 1.0 {
		return 0xFFFF
	}
	return uint32(v * 0xFFFF)
}

func U32ToF64(r, g, b uint32) (rf, gf, bf float64) { return u2f(r), u2f(g), u2f(b) }
func F64ToU32(rf, gf, bf float64) (r, g, b uint32) { return f2u(rf), f2u(gf), f2u(bf) }
