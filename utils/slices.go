package utils

func Filter[T any](slice []T, f func(T) bool) []T {
	var n = make([]T, 0)

	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}

	return n
}
