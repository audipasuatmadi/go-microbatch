package microbatch

type FlushStrategy[T any] interface {
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
