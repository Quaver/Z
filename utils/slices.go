package utils

// Filter Filters a slice based on a given condition
func Filter[T any](slice []T, f func(T) bool) []T {
	var n = make([]T, 0)

	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}

	return n
}

// Includes Returns if a slice contains an element
func Includes[T comparable](slice []T, val T) bool {
	for _, x := range slice {
		if x == val {
			return true
		}
	}

	return false
}
