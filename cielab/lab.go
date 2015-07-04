// cielab implements some CIELab color space conversions with D65 white point
package cielab

import "math"

const (
	e = 0.008856
	k = 903.3
)

var e_3 = math.Cbrt(e)

func toF(xr float64) (fx float64) {
	if xr > e {
		return math.Cbrt(xr)
	}
	return (k*xr + 16) / 116
}

func XYZToLAB(x, y, z float64) (l, a, b float64) {
	fx := toF(x / 0.95047)
	fy := toF(y / 1.00000)
	fz := toF(z / 1.08883)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)
	return
}

func LABToXYZ(l, a, b float64) (x, y, z float64) {
	fy := (l + 16) / 116
	fx := a/500 + fy
	fz := fy - b/200

	var xr, yr, zr float64
	if fx > e_3 {
		xr = math.Pow(fx, 3)
	} else {
		xr = (116*fx - 16) / k
	}

	if l > k*e {
		yr = math.Pow(fy, 3)
	} else {
		yr = l / k
	}

	if fz > e_3 {
		zr = math.Pow(fz, 3)
	} else {
		zr = (116*fz - 16) / k
	}

	return xr * 0.95047, yr * 1.00000, zr * 1.08883
}
