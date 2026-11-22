package eventually

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrTimeout = errors.New("eventually timed out")

// Run calls fn repeatedly, waiting interval between attempts, until
// fn returns nil or timeout elapses. fn receives a context that is cancelled
// when the timeout is reached, so long-running attempts can stop early.
//
// On timeout it returns an error wrapping both ErrTimeout and the last error
// returned by fn. If the parent ctx is cancelled first, ctx.Err() is wrapped.
func Run(
	ctx context.Context,
	timeout time.Duration,
	interval time.Duration,
	fn func(ctx context.Context) error,
) error {
	if interval <= 0 {
		return errors.New("eventually: interval must be > 0")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var lastErr error
	attempts := 0
	for {
		attempts++
		if lastErr = fn(ctx); lastErr == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return fmt.Errorf("%w after %s (%d attempts): %w", ErrTimeout, timeout, attempts, lastErr)
			}
			return fmt.Errorf("eventually: %w (last error: %w)", ctx.Err(), lastErr)
		case <-ticker.C:
		}
	}
}
