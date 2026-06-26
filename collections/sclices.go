package collections

import (
	"cmp"
	"slices"
)

func Merge[T cmp.Ordered](sets ...[]T) []T {
	result := []T{}
	for _, set := range sets {
		result = append(result, set...)
	}
	slices.Sort(result)
	return slices.Compact(result)
}
