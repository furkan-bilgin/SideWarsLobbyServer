package utils

import "math"

// https://stackoverflow.com/questions/15323767/does-go-have-if-x-in-construct-similar-to-python
func Contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

// Lerp function
func LinearInterp(v0 float64, v1 float64, t float64) float64 {
	t = math.Min(1, t)
	return v0*(1-t) + v1*t
}
