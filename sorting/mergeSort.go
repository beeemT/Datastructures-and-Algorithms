package sorting

func MergeSort(sort []int) []int {
	if len(sort) <= 1 {
		return sort
	}

	lS := len(sort) / 2
	sortedL := make([]int, lS)
	sortedR := make([]int, len(sort)-lS)

	if len(sort) > 2 {
		retChan := make(chan []int)

		go mergeSortChannel(sort[lS:], retChan)
		copy(sortedL, MergeSort(sort[:lS]))
		copy(sortedR, <-retChan)
	} else {
		sortedL[0] = sort[0]
		sortedR[0] = sort[1]
	}

	var iL, iR int
	lR := len(sortedR)
	lL := len(sortedL)
	for i := range sort {
		if (iL < lL) && (!(iR < lR) || (sortedL[iL] <= sortedR[iR])) {
			sort[i] = sortedL[iL]
			iL++
		} else {
			sort[i] = sortedR[iR]
			iR++
		}
	}

	return sort
}

func mergeSortChannel(sort []int, retChan chan []int) {
	retChan <- MergeSort(sort)
}
