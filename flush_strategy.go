package microbatch

type FlushStrategy[T any] interface {
	ShouldFlush(currentBatch Batch[T]) bool
}

type SizeBasedStrategy[T any] struct {
	MaxSize int
}

func (s *SizeBasedStrategy[T]) ShouldFlush(currentBatch Batch[T]) bool {
	return len(currentBatch) >= s.MaxSize
}
