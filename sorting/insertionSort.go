package sorting

func InsertionSort(a []int) {
	l := len(a)
	if l <= 1 {
		return
	}

	for i := 1; i < l; i++ {
		if a[i] >= a[i-1] {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			if a[j+1] < a[j] {
				a[j], a[j+1] = a[j+1], a[j]
				continue
			}
			break
		}
	}
}
