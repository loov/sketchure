package rgb

// r, g, b must be in range [0, 0xFFFF]
func LinearToXYZ(r, g, b uint32) (x, y, z float64) {
	rf, gf, bf := U32ToF64(r, g, b)
	x = 0.4124564*rf + 0.3575761*gf + 0.1804375*bf
	y = 0.2126729*rf + 0.7151522*gf + 0.0721750*bf
	z = 0.0193339*rf + 0.1191920*gf + 0.9503041*bf
	return
}

// r, g, b will be in range [0, 0xFFFF]
func LinearFromXYZ(x, y, z float64) (r, g, b uint32) {
	rf := 3.2404542*x - 1.5371385*y - 0.4985314*z
	gf := -0.9692660*x + 1.8760108*y + 0.0415560*z
	bf := 0.0556434*x - 0.2040259*y + 1.0572252*z
	return F64ToU32(rf, gf, bf)
}

// r, g, b must be in range [0, 0xFFFF]
func ToXYZ(r, g, b uint32) (x, y, z float64) {
	return LinearToXYZ(LinearizeU32(r, g, b))
}

// r, g, b will be in range [0, 0xFFFF]
func FromXYZ(x, y, z float64) (r, g, b uint32) {
	return DelinearizeU32(LinearFromXYZ(x, y, z))
}
