package eventbus

import "sync"

type Publisher[T any] interface {
	Publish(event T)
	TryPublish(event T)
}

type Subscriber[T any] interface {
	Subscribe() <-chan T
	Unsubscribe(ch <-chan T) bool
}

type EventBus[T any] struct {
	mu          sync.RWMutex
	subscribers map[<-chan T]chan T
	subDones    map[<-chan T]chan struct{}
	closed      bool
	done        chan struct{}
	wg          sync.WaitGroup
	bufferSize  int
	onDrop      func(T)
}

func NewEventBus[T any](bufferSize int, opts ...EventBusOption[T]) *EventBus[T] {
	if bufferSize < 0 {
		panic("eventbus: buffer size must be non-negative")
	}

	b := &EventBus[T]{
		subscribers: make(map[<-chan T]chan T),
		subDones:    make(map[<-chan T]chan struct{}),
		bufferSize:  bufferSize,
		done:        make(chan struct{}),
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *EventBus[T]) TryPublish(event T) {
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
	onDrop := b.onDrop
	b.mu.RUnlock()

	if onDrop == nil {
		return
	}
	for range dropped {
		onDrop(event)
	}
}

func (b *EventBus[T]) Publish(event T) {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return
	}

	b.wg.Add(1)
	pending := make(map[chan T]chan struct{}, len(b.subscribers))
	for _, ch := range b.subscribers {
		pending[ch] = b.subDones[ch]
	}
	done := b.done
	b.mu.RUnlock()
	defer b.wg.Done()

	for ch, subDone := range pending {
		select {
		case ch <- event:
		case <-subDone:
		case <-done:
			return
		}
	}
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
	b.subscribers[ch] = ch
	b.subDones[ch] = make(chan struct{})
	return ch
}

func (b *EventBus[T]) Unsubscribe(ch <-chan T) bool {
	if ch == nil {
		return false
	}

	b.mu.Lock()
	sendCh, ok := b.subscribers[ch]
	if ok {
		subDone := b.subDones[ch]
		delete(b.subscribers, ch)
		delete(b.subDones, ch)
		close(subDone)
	}
	b.mu.Unlock()

	if !ok {
		return false
	}

	b.wg.Wait()
	close(sendCh)
	return true
}

func (b *EventBus[T]) Close() {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return
	}

	b.closed = true
	close(b.done)
	subs := make([]chan T, 0, len(b.subscribers))
	for _, ch := range b.subscribers {
		subs = append(subs, ch)
	}
	b.subscribers = nil
	b.subDones = nil
	b.mu.Unlock()

	b.wg.Wait()
	for _, ch := range subs {
		close(ch)
	}
}

type EventBusOption[T any] func(*EventBus[T])

func WithDropHandler[T any](fn func(T)) EventBusOption[T] {
	return func(b *EventBus[T]) {
		b.onDrop = fn
	}
}
