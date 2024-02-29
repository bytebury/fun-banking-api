package utils

/**
 * Turns a nil value to a list or returns the given list.
 */
func Listify[T any](items []T) []T {
	if items == nil {
		return []T{}
	}
	return items
}
