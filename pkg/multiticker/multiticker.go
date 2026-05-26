package multiticker

import (
	"sync"
	"time"
)

const NoOffset time.Duration = -1

type MultiTicker struct {
	mu          sync.RWMutex
	ticker      *time.Ticker
	subscribers map[<-chan time.Time]chan time.Time
	stopCh      chan struct{}
	closed      bool
	interval    time.Duration
	offset      time.Duration
	onDrop      func(time.Time)
}

func NewMultiTicker(interval, offset time.Duration, opts ...MultiTickerOption) *MultiTicker {
	if interval <= 0 {
		panic("multiticker: interval must be positive")
	}
	if offset >= 0 {
		offset %= interval
	}

	t := &MultiTicker{
		subscribers: make(map[<-chan time.Time]chan time.Time),
		stopCh:      make(chan struct{}),
		interval:    interval,
		offset:      offset,
	}
	for _, opt := range opts {
		opt(t)
	}

	go func() {
		if offset >= 0 {
			var (
				intNs   = interval.Nanoseconds()
				offNs   = offset.Nanoseconds()
				elapsed = time.Now().UnixNano() % intNs
				sleep   time.Duration
			)
			if elapsed <= offNs {
				sleep = time.Duration(offNs - elapsed)
			} else {
				sleep = time.Duration((intNs - elapsed) + offNs)
			}
			if sleep > 0 {
				timer := time.NewTimer(sleep)
				select {
				case <-timer.C:
				case <-t.stopCh:
					timer.Stop()
					return
				}
			}
		}

		ticker := time.NewTicker(interval)
		t.mu.Lock()
		if t.closed {
			t.mu.Unlock()
			ticker.Stop()
			return
		}
		t.ticker = ticker

		dropped := 0
		now := time.Now()
		for _, ch := range t.subscribers {
			select {
			case ch <- now:
			default:
				dropped++
			}
		}
		onDrop := t.onDrop
		t.mu.Unlock()
		if onDrop != nil {
			for range dropped {
				onDrop(now)
			}
		}

		t.tick()
	}()

	return t
}

func (t *MultiTicker) Subscribe() <-chan time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		ch := make(chan time.Time)
		close(ch)
		return ch
	}

	ch := make(chan time.Time, 1)
	t.subscribers[ch] = ch
	return ch
}

func (t *MultiTicker) Unsubscribe(ch <-chan time.Time) bool {
	if ch == nil {
		return false
	}

	t.mu.Lock()
	sendCh, ok := t.subscribers[ch]
	if ok {
		delete(t.subscribers, ch)
	}
	t.mu.Unlock()

	if !ok {
		return false
	}
	close(sendCh)
	return true
}

func (t *MultiTicker) Stop() {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return
	}
	t.closed = true
	close(t.stopCh)
	if t.ticker != nil {
		t.ticker.Stop()
	}
	subs := make([]chan time.Time, 0, len(t.subscribers))
	for _, ch := range t.subscribers {
		subs = append(subs, ch)
	}
	t.subscribers = nil
	t.mu.Unlock()

	for _, ch := range subs {
		close(ch)
	}
}

func (t *MultiTicker) tick() {
	for {
		select {
		case tick := <-t.ticker.C:
			t.mu.RLock()
			if t.closed {
				t.mu.RUnlock()
				return
			}
			dropped := 0
			for _, ch := range t.subscribers {
				select {
				case ch <- tick:
				default:
					dropped++
				}
			}
			onDrop := t.onDrop
			t.mu.RUnlock()
			if onDrop != nil {
				for range dropped {
					onDrop(tick)
				}
			}
		case <-t.stopCh:
			return
		}
	}
}

type MultiTickerOption func(*MultiTicker)

func WithDropHandler(fn func(time.Time)) MultiTickerOption {
	return func(t *MultiTicker) {
		t.onDrop = fn
	}
}
