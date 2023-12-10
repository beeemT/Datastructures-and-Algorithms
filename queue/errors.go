package queue

import (
	"github.com/pkg/errors"
)

var (
	// ErrEmptyQueue is the error that is returned on operations that encounter an empty
	// queue but need a queue with elements inside.
	ErrEmptyQueue = errors.New("queue is empty")

	// ErrIndexOutOfBounds is the error that is returned on operations where an index is provided
	// but that index is not within the addressable space of the queue.
	ErrIndexOutOfBounds = errors.New("provided index is out of bounds")

	// ErrInvalidQueueType is returned when a nonexistent queuetype is encountered
	ErrInvalidQueueType = errors.New("provided queuetype is not invalid")

	// ErrInvalidQueueLimit is returned when a limit < 0 for the queue is encountered
	ErrInvalidQueueLimit = errors.New("provided limit for queue is invalid")
)
