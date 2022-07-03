package goutils

func Includes[T comparable](slice []T, value T) bool {
	for _, it := range slice {
		if it == value {
			return true
		}
	}
	return false
}
