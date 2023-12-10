package queue

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

// Iterator returns a channel which streams all elements of the queue.
// The amount of items cached in the channel can be determined by channelCapacity.
// The iterator can be stopped prematurely with the returned cancel function.
// Behaviour on concurrent calls to the queue iterator is undefined.
func (q *Queue[T]) Iterator(channelCapacity int) (<-chan T, context.CancelFunc) {
	ch := make(chan T, channelCapacity)
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context, cancel context.CancelFunc) {
		defer func() {
			if !errors.Is(ctx.Err(), context.Canceled) {
				cancel()
			}
			close(ch)
		}()

		for _, elem := range q.queueSlice {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- elem.Content()
			}
		}
	}(ctx, cancel)

	return ch, cancel
}

// MapInPlace executes the given mapping function on all elements in the queue in place.
func (q *Queue[T]) MapInPlace(f func(T) (T, error)) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.MapInPlaceUnsecure(f)
}

// MapInPlaceUnsecure executes the given mapping function on all elements in the queue in place.
func (q *Queue[T]) MapInPlaceUnsecure(f func(T) (T, error)) error {
	for i, elem := range q.queueSlice {
		newContent, err := f(elem.Content())
		if err != nil {
			return errors.Wrapf(err, "mapping element at position %d", i)
		}
		elem.SetContent(newContent)
	}

	return nil
}

// FilterInPlace executes the given filter function on all elements in the queue in place.
// Removes all elements for which the filter function returns false.
// Locks q.
func (q *Queue[T]) FilterInPlace(f func(T) (bool, error)) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.FilterInPlaceUnsecure(f)
}

// FilterInPlaceUnsecure executes the given filter function on all elements in the queue in place.
// Removes all elements for which the filter function returns false.
// Does not lock q.
func (q *Queue[T]) FilterInPlaceUnsecure(f func(T) (bool, error)) error {
	for i, elem := range q.queueSlice {
		if keep, err := f(elem.Content()); err == nil && !keep {
			_, err := q.remove(i)
			if err != nil {
				return errors.Wrapf(err, "filtering element at position %d", i)
			}
		} else if err != nil {
			return errors.Wrapf(err, "filtering element at position %d", i)
		}
	}

	return nil
}

// Fold executes a right fold fold function on all elements in the queue.
// Locks the queue.
func Fold[Aggregate, T any](
	q *Queue[T],
	initial Aggregate,
	f func(Aggregate, Element[T]) (Aggregate, error),
) (Aggregate, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	return FoldUnsecure(q, initial, f)
}

// FoldUnsecure executes a right fold fold function on all elements in the queue.
// Does not lock the queue.
func FoldUnsecure[Aggregate, T any](
	q *Queue[T],
	initial Aggregate,
	f func(Aggregate, Element[T]) (Aggregate, error),
) (Aggregate, error) {
	aggregate := initial
	for i, elem := range q.queueSlice {
		aggregate, err := f(aggregate, elem)
		if err != nil {
			return aggregate, errors.Wrapf(err, "folding element at position %d", i)
		}
	}

	return aggregate, nil
}

// Map clones the queue and executes the given mapping function on all elements in the new queue.
// The mapping function is responsible for the element projection and can determine whether the item
// should be included in the new queue.
// Locks q.
func Map[Told, Tnew any](
	q *Queue[Told],
	f func(Element[Told]) (Element[Tnew], bool, error),
) (*Queue[Tnew], error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	return MapUnsecure(q, f)
}

// MapUnsecure clones the queue and executes the given mapping function on all elements in the new
// queue.
// The mapping function is responsible for the element projection and can determine whether the item
// should be included in the new queue.
// Does not lock q.
func MapUnsecure[Told, Tnew any](
	q *Queue[Told],
	f func(Element[Told]) (Element[Tnew], bool, error),
) (*Queue[Tnew], error) {
	newQueue := &Queue[Tnew]{
		order:          q.order,
		queueSlice:     make([]Element[Tnew], q.numElements),
		numElements:    q.numElements,
		maxnumElements: q.maxnumElements,
		lock:           sync.Mutex{},
	}

	for i, elem := range q.queueSlice {
		if newElem, insert, err := f(elem); insert && err == nil {
			newQueue.Insert(newElem)
		} else if err != nil {
			return nil, errors.Wrapf(err, "mapping element at position %d", i)
		}
	}

	return newQueue, nil
}
