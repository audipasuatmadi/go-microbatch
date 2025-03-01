package microbatch

import "time"

type Event[T any] struct {
	Key     string
	Payload T
	AddedAt time.Time
	flush   bool
}

func (e *Event[T]) ShouldFlush() bool {
	return e.flush
}

func (e *Event[T]) Flush() {
	e.flush = true
}

type ResultEvent[T any] struct {
	Event Event[T]
	Err   error
}

type Batch[T any] []Event[T]

type ResultBatch[T any] []ResultEvent[T]
