package eventually_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/spacechunks/control-plane/internal/eventually"
)

func TestSucceedsEventually(t *testing.T) {
	calls := 0
	err := eventually.Run(context.Background(), time.Second, 10*time.Millisecond, func(context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("not yet")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestTimesOut(t *testing.T) {
	notReady := errors.New("not ready")
	start := time.Now()
	err := eventually.Run(context.Background(), 50*time.Millisecond, 10*time.Millisecond, func(context.Context) error {
		return notReady
	})
	if !errors.Is(err, eventually.ErrTimeout) || !errors.Is(err, notReady) {
		t.Fatalf("expected timeout wrapping last error, got %v", err)
	}
	if time.Since(start) > 200*time.Millisecond {
		t.Fatalf("took too long: %s", time.Since(start))
	}
}

func TestParentCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := eventually.Run(ctx, time.Second, 10*time.Millisecond, func(context.Context) error {
		return errors.New("nope")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
