package semantic

import (
	"context"
	"sync"
	"time"
)

// Five seconds between request starts stays below the account's 15 RPM limit.
// The shared limiter covers all Gemini instances within this process. Separate
// worker processes or other applications using the key need coordinated limits.
var providerRate = rateLimiter{interval: 5 * time.Second}

type rateLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	next     time.Time
}

func (r *rateLimiter) wait(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delay := time.Until(r.next)
	if delay > 0 {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	r.next = time.Now().Add(r.interval)
	return nil
}
