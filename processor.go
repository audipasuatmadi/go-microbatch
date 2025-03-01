package microbatch

type BatchProcessor[T any] interface {
	Process(data Batch[T]) ResultBatch[T]
}

// default batch processor
type simpleBatchProcessor[T any] struct{}

func (s *simpleBatchProcessor[T]) Process(data Batch[T]) ResultBatch[T] {
	results := make(ResultBatch[T], len(data))
	for i, d := range data {
		results[i] = ResultEvent[T]{Event: d, Err: nil}
	}

	return results
}
