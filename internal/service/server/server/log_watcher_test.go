package server

import (
	"context"
	"errors"
	"io"
	"regexp"
	"strings"
	"testing"
	"time"
)

// errReader returns its lines, then fails with err.
type errReader struct {
	r   io.Reader
	err error
}

func (e *errReader) Read(p []byte) (int, error) {
	n, err := e.r.Read(p)
	if err == io.EOF {
		return n, e.err
	}
	return n, err
}

func TestLogWatcherRun(t *testing.T) {
	readErr := errors.New("connection reset")

	tests := []struct {
		name    string
		reader  io.Reader
		pattern string
		err     error
	}{
		{
			name:    "match on first line",
			reader:  strings.NewReader("server started\n"),
			pattern: `server started`,
		},
		{
			name:    "match after other lines",
			reader:  strings.NewReader("booting\nloading config\nlistening on :8080\nmore output\n"),
			pattern: `listening on :\d+`,
		},
		{
			name:    "match on last line without trailing newline",
			reader:  strings.NewReader("booting\nready"),
			pattern: `^ready$`,
		},
		{
			name:    "pattern is matched per line, not across lines",
			reader:  strings.NewReader("server\nstarted\n"),
			pattern: `server\s+started`,
			err:     ErrNoMatch,
		},
		{
			name:    "stream ends without match",
			reader:  strings.NewReader("booting\nloading config\n"),
			pattern: `server started`,
			err:     ErrNoMatch,
		},
		{
			name:    "empty stream",
			reader:  strings.NewReader(""),
			pattern: `anything`,
			err:     ErrNoMatch,
		},
		{
			name:    "read error is returned",
			reader:  &errReader{r: strings.NewReader("booting\n"), err: readErr},
			pattern: `server started`,
			err:     readErr,
		},
		{
			name:    "match wins over later read error",
			reader:  &errReader{r: strings.NewReader("server started\n"), err: readErr},
			pattern: `server started`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			err := NewLogWatcher(tt.reader, regexp.MustCompile(tt.pattern)).Run(ctx)
			if !errors.Is(err, tt.err) {
				t.Fatalf("Run() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestLogWatcherMatchArrivesLater(t *testing.T) {
	pr, pw := io.Pipe()
	defer pr.Close()

	go func() {
		_, _ = io.WriteString(pw, "booting\n")
		time.Sleep(50 * time.Millisecond)
		_, _ = io.WriteString(pw, "server started\n")
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := NewLogWatcher(pr, regexp.MustCompile(`server started`)).Run(ctx); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
}

func TestLogWatcher_ContextCancelled(t *testing.T) {
	// a pipe nobody writes to blocks forever, like a quiet log stream.
	pr, pw := io.Pipe()
	defer pw.Close()
	defer pr.Close()

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(50*time.Millisecond, cancel)

	start := time.Now()
	err := NewLogWatcher(pr, regexp.MustCompile(`never`)).Run(ctx)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want %v", err, context.Canceled)
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("Run() took %s to return after cancel", elapsed)
	}
}

func TestLogWatcher_ContextTimeout(t *testing.T) {
	pr, pw := io.Pipe()
	defer pw.Close()
	defer pr.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := NewLogWatcher(pr, regexp.MustCompile(`never`)).Run(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Run() error = %v, want %v", err, context.DeadlineExceeded)
	}
}

func TestLogWatcher_AlreadyCancelled(t *testing.T) {
	pr, pw := io.Pipe()
	defer pw.Close()
	defer pr.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := NewLogWatcher(pr, regexp.MustCompile(`never`)).Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want %v", err, context.Canceled)
	}
}
