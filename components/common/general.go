package common

func RemoveElement[T comparable](target T, array []T) []T {
	for i, e := range array {
		if e == target {
			return append(array[:i], array[i+1:]...)
		}
	}
	return array
}

func HasElement[T comparable](target T, array []T) bool {
	for _, e := range array {
		if e == target {
			return true
		}
	}
	return false
}
