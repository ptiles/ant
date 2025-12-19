package geom

import (
	"math"
)

func FromDeg(deg float64) float64 {
	return deg / 180 * math.Pi
}

func Sin(deg float64) float64 {
	return math.Sin(FromDeg(deg))
}

func Cos(deg float64) float64 {
	return math.Cos(FromDeg(deg))
}

func SinOverSin(degA, degB float64) float64 {
	return Sin(degA) / Sin(degB)
}

// /    math.Phi   =  1.618...
// /    math.Phi-1 =  0.618...
// /  1-math.Phi   = -0.618...
// /   -math.Phi   = -1.618...
var sinOverSin5Exact = [4][4]float64{
	{+1.000000, math.Phi - 1, 1 - math.Phi, -1.000000},
	{+math.Phi, +1.000000000, -1.000000000, -math.Phi},
	{-math.Phi, -1.000000000, +1.000000000, +math.Phi},
	{-1.000000, 1 - math.Phi, math.Phi - 1, +1.000000},
}

// SinOverSin5 returns more precise values of SinOverSin() for GridLinesTotal == 5
func SinOverSin5(degA, degB float64) float64 {
	a := int(degA) / 72
	if a < 0 {
		a += 4
	} else {
		a -= 1
	}

	b := int(degB) / 72
	if b < 0 {
		b += 4
	} else {
		b -= 1
	}

	return sinOverSin5Exact[b][a]
}
