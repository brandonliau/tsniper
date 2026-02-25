package multiticker

import (
	"sync"
	"time"
)

const NoOffset time.Duration = -1

type MultiTicker struct {
	mu          sync.Mutex
	ticker      *time.Ticker
	subscribers []chan time.Time
	stopCh      chan struct{}
	stopped     bool
	onDrop      func()
}

func NewMultiTicker(interval, offset time.Duration, opts ...MultiTickerOption) *MultiTicker {
	if offset >= 0 && offset >= interval {
		offset = offset % interval
	}

	t := &MultiTicker{
		stopCh: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(t)
	}

	go func() {
		if offset >= 0 {
			var (
				now   = time.Now().UnixNano()
				intNs = interval.Nanoseconds()
				offNs = offset.Nanoseconds()
			)
			elapsed := now % intNs
			var sleep time.Duration
			if elapsed <= offNs {
				sleep = time.Duration(offNs - elapsed)
			} else {
				sleep = time.Duration((intNs - elapsed) + offNs)
			}
			timer := time.NewTimer(sleep)
			select {
			case <-timer.C:
			case <-t.stopCh:
				if !timer.Stop() {
					<-timer.C
				}
				return
			}
		}

		t.mu.Lock()
		if t.stopped {
			t.mu.Unlock()
			return
		}
		t.ticker = time.NewTicker(interval)

		tick := time.Now()
		for _, ch := range t.subscribers {
			select {
			case ch <- tick:
			default:
			}
		}
		t.mu.Unlock()

		t.tick()
	}()

	return t
}

func (t *MultiTicker) Subscribe() <-chan time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped {
		ch := make(chan time.Time)
		close(ch)
		return ch
	}

	c := make(chan time.Time, 1)
	t.subscribers = append(t.subscribers, c)
	return c
}

func (t *MultiTicker) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped {
		return
	}
	t.stopped = true
	close(t.stopCh)
	if t.ticker != nil {
		t.ticker.Stop()
	}
	for _, ch := range t.subscribers {
		close(ch)
	}
	t.subscribers = nil
}

func (t *MultiTicker) tick() {
	for {
		select {
		case tick := <-t.ticker.C:
			dropped := 0
			t.mu.Lock()
			for _, ch := range t.subscribers {
				select {
				case ch <- tick:
				default:
					dropped++
				}
			}
			t.mu.Unlock()
			if t.onDrop != nil {
				for range dropped {
					t.onDrop()
				}
			}
		case <-t.stopCh:
			return
		}
	}
}

type MultiTickerOption func(*MultiTicker)

func WithDropHandler(fn func()) MultiTickerOption {
	return func(t *MultiTicker) {
		t.onDrop = fn
	}
}
