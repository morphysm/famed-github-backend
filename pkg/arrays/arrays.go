package arrays

// Removes slice element at index(s) and returns new slice.
func Remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
