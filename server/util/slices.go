package util

import (
	"sort"
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

// SortMapByValue sorts a map by its values and returns the keys in descending order
func SortMapByValue(m map[int64]int64) []int64 {
	type kv struct {
		Key   int64
		Value int64
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	// Sort by value
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	var keys []int64
	for _, kv := range ss {
		keys = append(keys, kv.Key)
	}

	return keys
}

// SetDiff returns the set difference between two slices
func SetDiff[T comparable](a, b []T) []T {
	m := make(map[T]bool)
	for _, item := range b {
		m[item] = true
	}

	var diff []T
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}

	return diff
}
