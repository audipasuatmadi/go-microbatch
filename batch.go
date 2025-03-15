package microbatch

import (
	"time"
)

type Event[T any] struct {
	Payload T
	addedAt time.Time
}

func (e *Event[T]) AddedAt() time.Time {
	return e.addedAt
}

type Batch[T any] []Event[T]
