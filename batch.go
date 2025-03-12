package microbatch

import (
	"time"
)

type Event[T any] struct {
	Payload T
	addedAt time.Time
	flush   bool
}

func (e *Event[T]) ShouldFlush() bool {
	return e.flush
}

func (e *Event[T]) Flush() {
	e.flush = true
}

func (e *Event[T]) AddedAt() time.Time {
	return e.addedAt
}

type ResultEvent[T any] struct {
	Event Event[T]
	Err   error
}

type Batch[T any] []Event[T]

type ResultBatch[T any] []ResultEvent[T]
