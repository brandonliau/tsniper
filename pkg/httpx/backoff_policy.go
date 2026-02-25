package httpx

import (
	"math"
	"math/rand"
	"time"
)

type BackoffPolicy func(attempt int) time.Duration

func ConstantBackoff(delay time.Duration) BackoffPolicy {
	return func(attempt int) time.Duration {
		return delay
	}
}

func ExponentialBackoff(minDelay time.Duration, maxDelay time.Duration) BackoffPolicy {
	return func(attempt int) time.Duration {
		min := float64(minDelay.Milliseconds())
		max := float64(maxDelay.Milliseconds())
		delay := math.Min(min*math.Pow(2, float64(attempt)), max)
		return time.Duration(delay) * time.Millisecond
	}
}

func ExponentialBackoffWithJitter(minDelay time.Duration, maxDelay time.Duration) BackoffPolicy {
	return func(attempt int) time.Duration {
		min := float64(minDelay.Milliseconds())
		max := float64(maxDelay.Milliseconds())
		delay := math.Min(min*math.Pow(2, float64(attempt)), max)
		jitter := rand.Float64() * min
		return time.Duration(delay+jitter) * time.Millisecond
	}
}

func LinearBackoff(baseDelay time.Duration, maxDelay time.Duration) BackoffPolicy {
	return func(attempt int) time.Duration {
		delay := baseDelay * time.Duration(attempt+1)
		if delay > maxDelay {
			return maxDelay
		}
		return delay
	}
}

func NoBackoff() BackoffPolicy {
	return func(attempt int) time.Duration {
		return 0
	}
}
