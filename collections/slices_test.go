package collections

import "testing"

func TestMerge(t *testing.T) {
	a1 := []string{"a", "b", "c"}
	a2 := []string{"b", "c", "d", "f"}
	a3 := []string{"b", "d", "e"}
	result := Merge(a1, a2, a3)
	t.Logf("%v\n", result)
}
