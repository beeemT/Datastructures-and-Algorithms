package queue

// PeekElem returns a copy of the elem that would be returned on a call to Remove().
// Returns an error of type ErrEmptyQueue when the list is empty.
func (q *Queue[T]) PeekElem() (float64, T, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.numElements == 0 {
		return 0, *new(T), ErrEmptyQueue
	}
	elem := q.queueSlice[q.numElements-1] // dereference is a copy
	return elem.Priority(), elem.Content(), nil
}

// PeekElemAtIndex returns a copy of the elem at index.
// Returns an error of type ErrEmptyQueue when the list is empty.
// Returns an error of type ErrIndexOutOfBounds when the provided index is out of bounds.
func (q *Queue[T]) PeekElemAtIndex(index int) (float64, T, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.numElements == 0 {
		return 0, *new(T), ErrEmptyQueue
	}

	realIndex := (q.numElements - 1) - index
	if realIndex < 0 {
		return 0, *new(T), ErrIndexOutOfBounds
	}

	elem := q.queueSlice[realIndex] // dereference is a copy
	return elem.Priority(), elem.Content(), nil
}
