package microbatch

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Microbatch[T any] struct {
	wg sync.WaitGroup
	mu sync.Mutex

	ctx    Context
	stop   chan struct{}
	ticker *time.Ticker
	isOpen bool

	eventStream  chan Event[T]
	ResultStream chan ResultBatch[T]

	strategy  BatchStrategy[T]
	processor BatchProcessor[T]
}

func (m *Microbatch[T]) Add(ctx context.Context, events ...Event[T]) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isOpen {
		return ErrCantAddJob
	}

	for _, event := range events {
		m.eventStream <- event
	}

	return nil
}

func (m *Microbatch[T]) Start() {
	m.wg.Add(1)
	go m.run()
}

func (m *Microbatch[T]) run() {
	defer m.wg.Done()
	batch := Batch[T]{}

	for {
		select {
		case event, ok := <-m.eventStream:
			if !ok {
				return
			}

			batch = append(batch, event)
			if m.strategy.ShouldFlush(batch) {
				flushBatch := m.strategy.FlushBatch(batch)
				m.processBatch(flushBatch)

				if len(batch) == len(flushBatch) {
					batch = Batch[T]{}
				} else {
					remainingBatch := Batch[T]{}
					for _, e := range batch {
						if !m.strategy.ShouldFlush(Batch[T]{e}) {
							remainingBatch = append(remainingBatch, e)
						}
					}
					batch = remainingBatch
				}
			}

		case <-m.ticker.C:
			if len(batch) == 0 {
				continue
			}

			flushBatch := m.strategy.FlushBatch(batch)
			m.processBatch(flushBatch)

			if len(batch) == len(flushBatch) {
				batch = Batch[T]{}
			} else {
				remainingBatch := Batch[T]{}
				for _, e := range batch {
					if !m.strategy.ShouldFlush(Batch[T]{e}) {
						remainingBatch = append(remainingBatch, e)
					}
				}
				batch = remainingBatch
			}

		case <-m.stop:
			m.mu.Lock()
			m.isOpen = false
			m.mu.Unlock()
			close(m.eventStream)
			m.ticker.Stop()
			return
		}
	}
}

func (m *Microbatch[T]) processBatch(batch Batch[T]) {
	processed := m.processor.Process(batch)
	m.ResultStream <- processed
}

func (m *Microbatch[T]) Stop() {
	close(m.stop)
	m.wg.Wait()
}

type Config[T any] struct {
	Processor   BatchProcessor[T]
	Strategy    BatchStrategy[T]
	FlushTicker *time.Ticker
}

func New[T any](ctx Context, p Config[T]) (*Microbatch[T], error) {
	processor := p.Processor
	if processor == nil {
		processor = &simpleBatchProcessor[T]{}
	}

	strategy := p.Strategy
	if strategy == nil {
		return nil, errors.New("strategy is required")
	}

	ticker := p.FlushTicker
	if ticker == nil {
		ticker = time.NewTicker(5 * time.Second)
	}

	m := &Microbatch[T]{
		ctx:          ctx,
		stop:         make(chan struct{}),
		ticker:       ticker,
		isOpen:       true,
		eventStream:  make(chan Event[T]),
		ResultStream: make(chan ResultBatch[T]),
		strategy:     strategy,
		processor:    processor,
	}

	return m, nil
}
