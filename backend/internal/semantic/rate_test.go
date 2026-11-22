package semantic

import (
	"context"
	"testing"
	"time"
)

func TestRateWaitCanBeCancelled(t *testing.T) {
	r := rateLimiter{interval: time.Hour}
	if err := r.wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := r.wait(ctx); err != context.Canceled {
		t.Fatal("rate wait ignored cancellation", err)
	}
}

func TestRateSpacing(t *testing.T) {
	r := rateLimiter{interval: 5 * time.Millisecond}
	if err := r.wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	started := time.Now()
	if err := r.wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if time.Since(started) < 4*time.Millisecond {
		t.Fatal("requests were not spaced")
	}
}
