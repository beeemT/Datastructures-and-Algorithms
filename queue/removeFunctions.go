package queue

import (
	"math"

	"github.com/pkg/errors"
)

func (q *Queue[T]) remove(i int) (Element[T], error) {
	elem, err := q.deleteWithoutMemoryManagement(i)
	q.handleShrink()
	return elem, errors.Wrap(err, "removing element")
}

func (q *Queue[T]) handleShrink() {
	lenQ := len(q.queueSlice)
	if float64(lenQ) < q.shrinkFactor()*float64(cap(q.queueSlice)) {
		newCap := int(math.Ceil(q.afterShrinkFactor() * float64(cap(q.queueSlice))))
		temp := make([]Element[T], lenQ, newCap)
		copy(temp, q.queueSlice[:lenQ])
		q.queueSlice = temp
	}
}

func (q *Queue[T]) deleteWithoutMemoryManagement(i int) (Element[T], error) {
	lenQ := q.numElements
	if lenQ == 0 {
		return nil, ErrEmptyQueue
	}
	if i < 0 || i >= q.numElements {
		return nil, ErrIndexOutOfBounds
	}

	elem := q.queueSlice[i]
	if i == lenQ {
		q.queueSlice = q.queueSlice[:i]
	} else if i == 0 {
		q.queueSlice[0] = nil
		q.queueSlice = q.queueSlice[1:]
	} else {
		copy(q.queueSlice[i:], q.queueSlice[i+1:])
		q.queueSlice[lenQ-1] = nil
		q.queueSlice = q.queueSlice[:lenQ-1]
	}
	q.numElements--

	return elem, nil
}
