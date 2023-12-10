package queue

// PriorityElement encapsulates all information that is needed for the storage in the queue.
type PriorityElement[T any] struct {
	priority float64
	BaseElement[T]
}

func (e PriorityElement[T]) Priority() float64 {
	return e.priority
}

func (e *PriorityElement[T]) SetPriority(priority float64) {
	e.priority = priority
}

func (e PriorityElement[T]) Content() T {
	return *e.content
}

func (e *PriorityElement[T]) SetContent(content T) {
	e.content = &content
}

// BaseElement encapsulates all information that is needed for the storage in the queue.
type BaseElement[T any] struct {
	content *T
}

func (e BaseElement[T]) Priority() float64 {
	return 0
}

func (e *BaseElement[T]) SetPriority(priority float64) {
}

func (e BaseElement[T]) Content() T {
	return *e.content
}

func (e *BaseElement[T]) SetContent(content T) {
	e.content = &content
}
