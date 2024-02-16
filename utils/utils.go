package utils

import "math"

func FromDegrees(deg int) float64 {
	return float64(deg) * math.Pi / 180.0
}
