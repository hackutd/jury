package util

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

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

// IndexFunc from golang slices library
func IndexFunc[S ~[]E, E any](s S, f func(E) bool) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

// ContainsFunc from golang slices library
func ContainsFunc[S ~[]E, E any](s S, f func(E) bool) bool {
	return IndexFunc(s, f) >= 0
}

// Converts an int slice to string slice
func IntToString[T constraints.Integer](arr []T) []string {
	out := make([]string, len(arr))
	for i, v := range arr {
		out[i] = strconv.Itoa(int(v))
	}
	return out
}
