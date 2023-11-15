package util

// Map applies a function to each element of a slice and returns a new slice
func Map[T, U interface{}](arr []T, fn func(T) U) []U {
	out := make([]U, len(arr))
	for i, v := range arr {
		out[i] = fn(v)
	}
	return out
}

// Any returns true if any element of the slice is true
func Any[T bool](arr []T) bool {
	for _, v := range arr {
		if v {
			return true
		}
	}
	return false
}

// All returns true if all elements of the slice are true
func All[T bool](arr []T) bool {
	for _, v := range arr {
		if !v {
			return false
		}
	}
	return true
}
