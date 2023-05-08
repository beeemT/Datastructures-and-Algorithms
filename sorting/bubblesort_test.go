package sorting

import (
	"sort"
	"testing"
)

func TestBubbleSortIntSlice(t *testing.T) {
	t.Parallel()
	data := make([]int, len(ints))
	copy(data, ints)
	BubbleSort(data)
	if !sort.IsSorted(sort.IntSlice(data)) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}
