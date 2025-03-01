package microbatch

import "time"

type BatchStrategy[T any] interface {
	ShouldFlush(currentBatch Batch[T]) bool
	FlushBatch(currentBatch Batch[T]) Batch[T]
}

type SizeBasedStrategy[T any] struct {
	MaxSize int
}

func (s *SizeBasedStrategy[T]) ShouldFlush(currentBatch Batch[T]) bool {
	return len(currentBatch) >= s.MaxSize
}

func (s *SizeBasedStrategy[T]) FlushBatch(currentBatch Batch[T]) Batch[T] {
	if len(currentBatch) == 0 {
		return nil
	}

	if len(currentBatch) <= s.MaxSize {
		return currentBatch
	}

	toFlushEvents := Batch[T]{}
	for i := s.MaxSize; i < len(currentBatch); i++ {
		currentBatch[i].Flush()
		toFlushEvents = append(toFlushEvents, currentBatch[i])
	}

	return toFlushEvents
}

type TimeBasedStrategy[T any] struct {
	flushedAt        *time.Time
	IntervalDuration time.Duration
}

func (t *TimeBasedStrategy[T]) ShouldFlush(events Batch[T]) bool {
	if t.flushedAt == nil {
		return true
	}

	for _, event := range events {
		if event.AddedAt.After(*t.flushedAt) {
			return true
		}
	}

	return false
}

func (t *TimeBasedStrategy[T]) FlushBatch(events Batch[T]) Batch[T] {
	if t.flushedAt == nil {
		initFlushedAt := time.Now()
		t.flushedAt = &initFlushedAt
	}

	if len(events) == 0 {
		return nil
	}

	toFlushEvents := Batch[T]{}
	for _, event := range events {
		if event.AddedAt.After(*t.flushedAt) {
			event.Flush()
			toFlushEvents = append(toFlushEvents, event)
		}
	}

	return toFlushEvents
}
