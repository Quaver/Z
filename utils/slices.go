package utils

import "errors"

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

// Find Finds an element in a slice according to a given condition
func Find[T comparable](slice []T, f func(T) bool) (T, error) {
	for _, e := range slice {
		if f(e) {
			return e, nil
		}
	}

	return *new(T), errors.New("could not find element in slice")
}

// FindIndex Returns the index of an element in a slice
func FindIndex[T comparable](slice []T, val T) int {
	for i, x := range slice {
		if x == val {
			return i
		}
	}

	return -1
}
