package microbatch

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Microbatch[T any] struct {
	wg sync.WaitGroup
	mu sync.Mutex

	ctx    context.Context
	stop   chan struct{}
	isOpen bool

	batchCtx             context.Context
	batchCancelCtx       context.CancelFunc
	batchTimeoutDuration time.Duration

	eventStream  chan Event[T]
	ResultStream chan ResultBatch[T]

	strategy  FlushStrategy[T]
	processor BatchProcessor[T]
}

func (m *Microbatch[T]) Add(ctx context.Context, events ...T) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isOpen {
		return fmt.Errorf("%w: microbatch is closed", ErrCantAddJob)
	}

	for _, event := range events {
		m.eventStream <- Event[T]{Payload: event, addedAt: time.Now()}
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

	flush := func() {
		flushBatch := m.strategy.FlushBatch(batch)
		m.processBatch(flushBatch)

		if len(batch) == len(flushBatch) {
			batch = Batch[T]{}
		} else {
			remainingBatch := Batch[T]{}
			for _, e := range batch {
				if !e.ShouldFlush() {
					remainingBatch = append(remainingBatch, e)
				}
			}
			batch = remainingBatch
		}
	}

	for {
		select {
		case event, ok := <-m.eventStream:
			if !ok {
				return
			}
			batch = append(batch, event)
			if !m.strategy.ShouldFlush(batch) {
				continue
			}
			flush()

		// TODO: implement context cancellation from m.ctx
		case <-m.batchCtx.Done():
			if len(batch) == 0 {
				continue
			}
			flush()
			m.batchCtx, m.batchCancelCtx = context.WithTimeout(m.ctx, m.batchTimeoutDuration)

		case <-m.stop:
			m.mu.Lock()
			m.isOpen = false
			m.mu.Unlock()
			close(m.eventStream)
			return
		}
	}
}

func (m *Microbatch[T]) processBatch(batch Batch[T]) {
	m.ResultStream <- m.processor.Process(batch)
}

func (m *Microbatch[T]) Stop() {
	close(m.stop)
	m.wg.Wait()
}

type Config[T any] struct {
	Processor    BatchProcessor[T]
	Strategy     FlushStrategy[T]
	BatchTimeout *time.Duration
}

func New[T any](ctx context.Context, p Config[T]) (*Microbatch[T], error) {
	processor := p.Processor
	if processor == nil {
		processor = &simpleBatchProcessor[T]{}
	}

	strategy := p.Strategy
	if strategy == nil {
		strategy = &SizeBasedStrategy[T]{MaxSize: 5}
	}

	batchTimeout := defaultBatchTimeoutDuration
	if p.BatchTimeout != nil {
		batchTimeout = *p.BatchTimeout
	}

	batchCtx, batchCtxCancel := context.WithTimeout(ctx, batchTimeout)

	return &Microbatch[T]{
		ctx:                  ctx,
		batchCtx:             batchCtx,
		batchCancelCtx:       batchCtxCancel,
		batchTimeoutDuration: batchTimeout,
		stop:                 make(chan struct{}),
		isOpen:               true,
		eventStream:          make(chan Event[T]),
		ResultStream:         make(chan ResultBatch[T]),
		strategy:             strategy,
		processor:            processor,
	}, nil
}
