package util

import "sort"

// SortIDList sorts an ID list.
func SortIDList[T ~int](list []T) {
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
}

// EqualIDList returns whether the two supplied ID lists are the same.  Both
// are assumed to be sorted.
func EqualIDList[T ~int](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
