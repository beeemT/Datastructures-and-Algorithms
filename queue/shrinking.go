package queue

// shrinkFactor determines a factor dynamically depending on the amount of elements in the queue
// at what point to initiate a shrink operation on the underlying slice
func (q *Queue[T]) shrinkFactor() float64 {
	switch {
	case q.numElements < 1000:
		return 0.75
	case q.numElements < 10000:
		return 0.9
	case q.numElements < 100000:
		return 0.99
	case q.numElements < 1000000:
		return 0.999
	default:
		return 0.9999
	}
}

// afterShrinkFactor determines a factor dynamically depending on the amount of elements in the
// queue by how much to shrink the underlying slice on a shrink operation
func (q *Queue[T]) afterShrinkFactor() float64 {
	switch {
	case q.numElements < 1000:
		return 0.8
	case q.numElements < 10000:
		return 0.95
	case q.numElements < 100000:
		return 0.995
	case q.numElements < 1000000:
		return 0.9995
	default:
		return 0.99995
	}
}
