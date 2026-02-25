package eventbus

import (
	"sync"
)

type Publisher[T any] interface {
	Publish(event T)
	PublishBlocking(event T)
}

type Subscriber[T any] interface {
	Subscribe() <-chan T
}

type EventBus[T any] struct {
	mu          sync.RWMutex
	subscribers []chan T
	closed      bool
	bufferSize  int
	onDrop      func(T)
}

func NewEventBus[T any](bufferSize int, opts ...EventBusOption[T]) *EventBus[T] {
	b := &EventBus[T]{
		bufferSize: bufferSize,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *EventBus[T]) Publish(event T) {
	b.mu.RLock()

	if b.closed {
		b.mu.RUnlock()
		return
	}

	dropped := 0
	for _, ch := range b.subscribers {
		select {
		case ch <- event:
		default:
			dropped++
		}
	}
	b.mu.RUnlock()

	if b.onDrop != nil {
		for range dropped {
			b.onDrop(event)
		}
	}
}

func (b *EventBus[T]) PublishBlocking(event T) {
	b.mu.RLock()

	if b.closed {
		b.mu.RUnlock()
		return
	}

	for _, ch := range b.subscribers {
		ch <- event
	}
	b.mu.RUnlock()
}

func (b *EventBus[T]) Subscribe() <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		ch := make(chan T)
		close(ch)
		return ch
	}

	ch := make(chan T, b.bufferSize)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *EventBus[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	b.closed = true
	for _, ch := range b.subscribers {
		close(ch)
	}
	b.subscribers = nil
}

type EventBusOption[T any] func(*EventBus[T])

func WithDropHandler[T any](fn func(T)) EventBusOption[T] {
	return func(b *EventBus[T]) {
		b.onDrop = fn
	}
}
