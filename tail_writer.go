package logger

import (
	"context"
	"regexp"
	"strings"
	"sync"
)

// tailWriter struct combines the Buffer interface with tailing capabilities.
type tailWriter struct {
	Buffer
	mu    sync.RWMutex
	max   uint
	lines []string
}

// Tail returns the collected log lines as a single string.
func (w *tailWriter) Tail() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return strings.Join(w.lines, "\n")
}

// ondata handles incoming buffer data and stores lines while maintaining a maximum number of lines.
func (w *tailWriter) ondata(buffer []byte) {
	reg := regexp.MustCompile(`[\r\n]+`)
	lines := reg.Split(strings.TrimRight(string(buffer), "\r\n"), -1)
	excessLines := len(w.lines) + len(lines) - int(w.max)
	if excessLines > 0 {
		newlines := make([]string, len(w.lines[excessLines:]))
		copy(newlines, w.lines[excessLines:])
		w.mu.Lock()
		w.lines = append(newlines, lines...)
		w.mu.Unlock()
	} else {
		w.mu.Lock()
		w.lines = append(w.lines, lines...)
		w.mu.Unlock()
	}
}

// TailWriter initializes and returns a new tailWriter instance with the given parameters.
func TailWriter(ctx context.Context, max uint, buffer int) *tailWriter {
	w := &tailWriter{
		max: max,
	}
	w.Buffer = NewBuffer(ctx, buffer, w.ondata, nil, nil)
	return w
}
