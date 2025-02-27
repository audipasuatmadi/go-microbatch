package microbatch

import (
	"context"
	"sync"
	"time"
)

type Microbatch[T any] interface {
	Add(ctx context.Context, data T)
	ReadData(ctx context.Context) (data []T)
}

type microbatch[T any] struct {
	mu sync.Mutex

	buffer         []T
	bufferIterator int32
	maxBufferSize  int32
	flushInterval  time.Duration
	fullChan       chan bool
	hasResetChan   chan bool
}

func (m *microbatch[T]) Add(ctx context.Context, data T) {
	m.mu.Lock()
	if m.bufferIterator < m.maxBufferSize {
		m.buffer[m.bufferIterator] = data
		m.bufferIterator++

		if m.bufferIterator == m.maxBufferSize {
			m.fullChan <- true
			<-m.hasResetChan
		}
	}
	m.mu.Unlock()

}

func (m *microbatch[T]) ReadData(ctx context.Context) (data []T) {
	<-m.fullChan
	var allData []T = make([]T, 0)
	for i := int32(0); i < m.bufferIterator; i++ {
		allData = append(allData, m.buffer[i])
	}

	m.buffer = make([]T, m.maxBufferSize)
	m.bufferIterator = 0
	m.hasResetChan <- true
	return allData
}

func (m *microbatch[T]) batchByInterval() {
	for {
		time.Sleep(m.flushInterval)
		if m.bufferIterator > 0 {
			m.mu.Lock()
			m.fullChan <- true
			<-m.hasResetChan
			m.mu.Unlock()
		}
	}
}

type NewParams struct {
	MaxSize       int32
	FlushInterval time.Duration
}

func New[T any](p NewParams) Microbatch[T] {
	m := &microbatch[T]{
		buffer:         make([]T, p.MaxSize),
		bufferIterator: 0,
		maxBufferSize:  p.MaxSize,
		flushInterval:  p.FlushInterval,
		fullChan:       make(chan bool, 1),
		hasResetChan:   make(chan bool, 1),
	}

	if p.FlushInterval != 0 {
		go m.batchByInterval()
	}

	return m
}
