package lab

import (
	"math"

	"github.com/loov/sketch-capture/rgb"
)

type WhitePoint struct{ X, Y, Z float64 }

// This is the default reference white point.
var (
	D65 = WhitePoint{X: 0.95047, Y: 1.00000, Z: 1.08883}
	D50 = WhitePoint{X: 0.96422, Y: 1.00000, Z: 0.82521}

	DefaultWhitePoint = D65
)

func FromXYZ(x, y, z float64) (l, a, b float64) { return FromXYZRef(x, y, z, DefaultWhitePoint) }
func ToXYZ(l, a, b float64) (x, y, z float64)   { return ToXYZRef(l, a, b, DefaultWhitePoint) }

const (
	e = 0.008856
	k = 903.3
)

var (
	e_3 = math.Cbrt(e)
)

func toF(xr float64) (fx float64) {
	if xr > e {
		return math.Cbrt(xr)
	}
	return (k*xr + 16) / 116
}

func FromXYZRef(x, y, z float64, white WhitePoint) (l, a, b float64) {
	fx := toF(x / white.X)
	fy := toF(y / white.Y)
	fz := toF(z / white.Z)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)
	return
}

func ToXYZRef(l, a, b float64, white WhitePoint) (x, y, z float64) {
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

	return xr * white.X, yr * white.Y, zr * white.Z
}

func FromRGB(r, g, b uint32) (l, la, lb float64) { return FromXYZ(rgb.ToXYZ(r, g, b)) }
func ToRGB(l, la, lb float64) (r, g, b uint32)   { return rgb.FromXYZ(ToXYZ(l, la, lb)) }

func FromRGBRef(r, g, b uint32, white WhitePoint) (l, la, lb float64) {
	x, y, z := rgb.ToXYZ(r, g, b)
	return FromXYZRef(x, y, z, white)
}
func ToRGBRef(l, la, lb float64, white WhitePoint) (r, g, b uint32) {
	return rgb.FromXYZ(ToXYZRef(l, la, lb, white))
}
