package sorting

import (
	"sort"
	"testing"
)

func TestMergeSortIntSlice(t *testing.T) {
	t.Parallel()
	data := make([]int, len(ints))
	copy(data, ints)
	ret := MergeSort(data)
	if !sort.IsSorted(sort.IntSlice(ret)) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", ret)
	}
}
