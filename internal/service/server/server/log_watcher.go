package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
)

var ErrNoMatch = fmt.Errorf("log stream ended without a match")

type LogWatcher struct {
	reader io.Reader
	regex  *regexp.Regexp
}

func NewLogWatcher(r io.Reader, regex *regexp.Regexp) *LogWatcher {
	return &LogWatcher{
		reader: r,
		regex:  regex,
	}
}

func (w *LogWatcher) Run(ctx context.Context) error {
	result := make(chan error)

	go func() {
		sc := bufio.NewScanner(w.reader)
		for sc.Scan() {
			if w.regex.Match(sc.Bytes()) {
				result <- nil
				return
			}
		}
		if err := sc.Err(); err != nil {
			result <- err
			return
		}
		result <- ErrNoMatch
	}()

	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
