package imagelib

import (
	"math"
)

// Smallest of three numbers
func min(a, b, c uint8) uint8 {
	if (a < b) && (a < c) {
		return a
	} else if (b < a) && (b < c) {
		return b
	}
	return c
}

// Largest of three numbers
func max(a, b, c uint8) uint8 {
	if (a >= b) && (a >= c) {
		return a
	} else if (b >= a) && (b >= c) {
		return b
	}
	return c
}

// Smallest of three floats
func fmin(a, b, c float64) float64 {
	return math.Min(math.Min(a, b), c)
}

// Largest of three floats
func fmax(a, b, c float64) float64 {
	return math.Max(math.Max(a, b), c)
}

// Absolute value
func abs(a int8) uint8 {
	if a < 0 {
		return uint8(-a)
	}
	return uint8(a)
}

// Absolute value
func fabs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
