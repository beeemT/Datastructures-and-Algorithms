package queue

import "github.com/pkg/errors"

func (q *Queue[T]) insertFifo(elem Element[T]) {
	q.queueSlice = append([]Element[T]{elem}, q.queueSlice...)
}

func (q *Queue[T]) insertLifo(elem Element[T]) {
	q.queueSlice = append(q.queueSlice, elem)
}

func (q *Queue[T]) insertPriorityHigh(elem Element[T]) {
	// If the queue is empty or the new element has a higher priority than the current item with the
	// highest priority
	// it can be appended to the slice.
	if q.numElements == 0 || (q.queueSlice[q.numElements-1]).Priority() < elem.Priority() {
		q.queueSlice = append(q.queueSlice, elem)
		return
	}

	if (q.queueSlice[q.numElements-1]).Priority() == elem.Priority() {
		q.backtrackInsertionPoint(elem)
	}

	// Default case. Iterate through full queue until the first suitable spot for the new element is
	// found.
	for i, e := range q.queueSlice {
		if e.Priority() < elem.Priority() {
			continue
		}

		// e.prio >= elem.prio
		q.queueSlice = append(
			q.queueSlice[:(i-1)],
			append([]Element[T]{elem}, q.queueSlice[(i-1):]...)...)
		break
	}
}

func (q *Queue[T]) insertPriorityLow(elem Element[T]) {
	// If the queue is empty or the new element has a lower priority than the current item with the
	// lowest priority
	// it can be appended to the slice.
	if q.numElements == 0 || (q.queueSlice[q.numElements-1]).Priority() > elem.Priority() {
		q.queueSlice = append(q.queueSlice, elem)
		return
	}

	if (q.queueSlice[q.numElements-1]).Priority() == elem.Priority() {
		q.backtrackInsertionPoint(elem)
	}

	// Default case. Iterate through full queue until the first suitable spot for the new element is
	// found.
	for i, e := range q.queueSlice {
		if e.Priority() > elem.Priority() {
			continue
		}

		// e.prio <= elem.prio
		q.queueSlice = append(
			q.queueSlice[:(i-1)],
			append([]Element[T]{elem}, q.queueSlice[(i-1):]...)...)
		break
	}
}

// backtrackInsertionPoint finds a suitable insertion point for an item from the back of the queue.
// queueSlice[q.numElements-1].Priority() == elem.Priority() doesn't need full iteration over queue
// to find the spot for insertion. Iteration from slice end to first element which has non equal
// priority than the new element is enough because of the priority invariant.
// In the worst case this will iterate over whole queue, so attach new element to front of slice.
func (q *Queue[T]) backtrackInsertionPoint(elem Element[T]) {
	for i := q.numElements - 1; i > -1; i-- {
		if q.queueSlice[i].Priority() == elem.Priority() {
			continue
		}
		q.queueSlice = append(q.queueSlice[i:], append([]Element[T]{elem}, q.queueSlice[:i]...)...)
		return
	}
	q.queueSlice = append([]Element[T]{elem}, q.queueSlice...)
}

func (q *Queue[T]) insertFifoLimited(elem Element[T]) error {
	if q.numElements == q.maxnumElements && q.maxnumElements != 0 {
		_, err := q.remove(q.numElements - 1)
		if err != nil {
			return errors.Wrap(err, "popping element because of overflow")
		}
	}
	q.insertFifo(elem)
	return nil
}
