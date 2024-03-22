package internal

func Filter[T any](array []T, predicate func(T) bool) []T {
	filtered := []T{}
	for _, v := range array {
		if predicate(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func Map[A any, B any](array []A, predicate func(A) B) []B {
	mapped := []B{}
	for _, v := range array {
		mapped = append(mapped, predicate(v))
	}

	return mapped
}
