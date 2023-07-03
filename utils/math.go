package utils

import "golang.org/x/exp/constraints"

func Clamp[T constraints.Integer | constraints.Float](f, low, high T) T {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}
