package queue

import (
	"sync"
)

// Queuetype is the enum type for queue invariants.
// Invariants:
//
//	`Always structure the slice in a way that the item at len(queueSlice)-1 is the item for the
//	remove operation
//
//	For same main ordering property of two elements the element that is older will be removed.
//	Fifo:
//		len(queueSlice)-1 is the fist inserted elem
//	Lifo:
//		len(queueSlice)-1 is the last inserted elem
//	PriorityHigh:
//		len(queueSlice)-1 is the elem with highest priority
//	PriorityLow:
//		len(queueSlice)-1 is the elem with lowest priority
type Queuetype int

const (
	// Fifo queue.
	Fifo Queuetype = iota

	// Lifo queue.
	Lifo

	// PriorityHigh means that on remove the elem with the highest priority value is returned.
	PriorityHigh

	// PriorityLow means that on remove the elem with the lowest priority value is returned.
	PriorityLow

	// FifoLimited means that the queue has a maximum capacity. Requires extra call to set capacity.
	FifoLimited

	numQueuetypes = 5
)

// Element is the interface encapsulating all element types
type Element[T any] interface {
	Priority() float64
	SetPriority(float64)

	Content() T
	SetContent(T)
}

// Queue is a queue of type Queuetype
type Queue[T any] struct {
	order          Queuetype
	lock           sync.Mutex
	queueSlice     []Element[T]
	numElements    int
	maxnumElements int
}

// NewQueue builds a new Queue with the passed Queuetype.
// Since the queue is realized through a slice, expectedLength is the initial
// cap() value of said slice.
func NewQueue[T any](tp Queuetype) (*Queue[T], error) {
	if tp < 0 || tp > numQueuetypes {
		return nil, ErrInvalidQueueType
	}

	return &Queue[T]{
		order:      tp,
		queueSlice: make([]Element[T], 0),
	}, nil
}

// NewPriorityElement builds a new Element with the passed content and priority.
// You cannot work with the element directly. This return value is only meant to be passed to
// queue functions.
func NewPriorityElement[T any](c T, priority float64) *PriorityElement[T] {
	return &PriorityElement[T]{
		priority:    priority,
		BaseElement: *NewBaseElement(c),
	}
}

// NewBaseElement builds a new Element with the passed content and priority = 0.
// You cannot work with the element directly. This return value is only meant to be passed to
// queue functions.
func NewBaseElement[T any](c T) *BaseElement[T] {
	return &BaseElement[T]{content: &c}
}

// Len returns the number of elements in the queue.
func (q *Queue[T]) Len() int {
	return q.numElements
}

// SetLimit sets the max capacity for the queue. Returns a ErrInvalidQueueLimit if limit < 0.
func (q *Queue[T]) SetLimit(limit int) error {
	if limit < 0 {
		return ErrInvalidQueueLimit
	}
	q.maxnumElements = limit
	return nil
}

// Append literally appends the element to the queue.
// Append does not uphold the invariant of the queue defined by the Queuetype and is thus unsafe.
// Use Insert for honoring the invariant.
func (q *Queue[T]) Append(elem Element[T]) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queueSlice = append(q.queueSlice, elem)
	q.numElements++
}

// Insert inserts the passed element into the queue, according to the Queuetype of the queue.
// Insert upholds the invariant of the Queue.
// When there are multiple elements with the same priority the oldest elem will be the first that is
// removed.
func (q *Queue[T]) Insert(elem Element[T]) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	switch q.order {
	case Fifo:
		q.insertFifo(elem)
	case Lifo:
		q.insertLifo(elem)
	case PriorityHigh:
		q.insertPriorityHigh(elem)
	case PriorityLow:
		q.insertPriorityLow(elem)
	case FifoLimited:
		return q.insertFifoLimited(elem)
	default:
		return ErrInvalidQueueType
	}
	q.numElements++
	return nil
}

// Remove pops the element that is meant to be removed first according to the queues order.
// When there are multiple elements with the same priority the oldest elem will be the first that is
// removed (FIFO).
// Returns the Element split up into its pieces.
// If the list is empty, an error is returned.
func (q *Queue[T]) Remove() (T, float64, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	elem, err := q.remove(q.numElements - 1)
	if err != nil {
		return *new(T), 0, err
	}
	return elem.Content(), elem.Priority(), nil
}

// RemoveElement pops the element that is meant to be removed first according to the queues order.
// When there are multiple elements with the same priority the oldest elem will be the first that is
// removed.
// Returns the pointer to the Element itself.
func (q *Queue[T]) RemoveElement() (Element[T], error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	elem, err := q.remove(q.numElements - 1)
	if err != nil {
		return nil, err
	}

	return elem, nil
}

// UpdatePriority updates the priority of all elements with priority oldPriority to the newPriority.
// Upholds the invariant of the queue.
// Returns the number of updates.
// If performanceFlag is set, elements with the same priority will be reversed in their order for
// ordertypes
// PriorityHigh and PriorityLow.
func (q *Queue[T]) UpdatePriority(oldPriority, newPriority float64, performanceFlag bool) int {
	q.lock.Lock()
	defer q.lock.Unlock()

	counter := 0

	var list []Element[T]
	if !performanceFlag {
		list = make([]Element[T], 0) // for buffering elements for reinsertion
	}

	switch q.order {
	case Lifo, Fifo:
		for _, e := range q.queueSlice { // O(n)
			//modifing e works because queueSlice is Element
			//+ Lifo and Fifo both are not sorted after priority

			if e.Priority() == oldPriority {
				e.SetPriority(newPriority)
				counter++
			}
		}

	case PriorityHigh, PriorityLow:
		// todo: use binsearch to find first elem with priority
		var modFlag bool

		for i, e := range q.queueSlice {
			if e.Priority() == oldPriority {
				q.deleteWithoutMemoryManagement(
					i,
				) // delete without MemoryManagement because elements get reinserted
				e.SetPriority(newPriority)
				if performanceFlag {
					q.Insert(e) // reverses the order within elements with the same priority
				} else {
					list = append(list, e)
				}
			} else if modFlag {
				break
			}
		}
	}

	if (q.order == PriorityHigh || q.order == PriorityLow) && !performanceFlag {
		l := len(list)
		for i := range list {
			q.Insert(list[l-(i+1)]) // insert oldest element first
		}
	}

	return counter
}

// GetAllElements returns a slice of all elements contents.
func (q *Queue[T]) GetAllElements() []T {
	ret := make([]T, q.numElements)
	for _, elem := range q.queueSlice {
		ret = append(ret, elem.Content())
	}
	return ret
}

// Clone clones the queue completely.
// Since only the elements can be realistically copied, if the element content is a reference type
// the original data in the queue can still be affected by changes on the new queue.
func (q *Queue[T]) Clone() *Queue[T] {
	q.lock.Lock()
	defer q.lock.Unlock()

	newQueue := &Queue[T]{
		order:          q.order,
		queueSlice:     make([]Element[T], q.numElements),
		numElements:    q.numElements,
		maxnumElements: q.maxnumElements,
		lock:           sync.Mutex{},
	}

	copy(newQueue.queueSlice, q.queueSlice)

	return newQueue
}
