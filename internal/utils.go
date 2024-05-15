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

func Map[A any, B any](array []A, predicate func(A) (B, error)) ([]B, error) {
	mapped := []B{}
	for _, v := range array {
		mappedValue, err := predicate(v)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, mappedValue)
	}

	return mapped, nil
}

func FlatMap[A any, B any](array []A, predicate func(A) ([]B, error)) ([]B, error) {
	mapped, err := Map(array, predicate)
	if err != nil {
		return nil, err
	}

	return Flatten(mapped), nil
}

func Flatten[A any](array [][]A) []A {
	flatmapped := []A{}
	for _, v := range array {
		flatmapped = append(flatmapped, v...)
	}
	return flatmapped
}
