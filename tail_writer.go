package logger

import (
	"context"
	"regexp"
	"strings"
	"sync"
)

// TailWriter struct combines the Buffer interface with tailing capabilities.
type TailWriter struct {
	Buffer
	mu    sync.RWMutex
	max   uint
	lines []string
}

// Tail returns the collected log lines as a single string.
func (w *TailWriter) Tail() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return strings.Join(w.lines, "\n")
}

// ondata handles incoming buffer data and stores lines while maintaining a maximum number of lines.
func (w *TailWriter) ondata(buffer []byte) {
	reg := regexp.MustCompile(`[\r\n]+`)
	lines := reg.Split(strings.TrimRight(string(buffer), "\r\n"), -1)
	w.mu.Lock()
	defer w.mu.Unlock()
	if len_lines, max := len(lines), int(w.max); len_lines > max {
		w.lines = lines[len_lines-max:]
	} else if excessLines := len(w.lines) + len_lines - max; excessLines > 0 {
		newlines := make([]string, max-len_lines)
		copy(newlines, w.lines[excessLines:])
		w.lines = append(newlines, lines...)
	} else {
		w.lines = append(w.lines, lines...)
	}
}

// NewTailWriter initializes and returns a new TailWriter instance with the given parameters.
func NewTailWriter(ctx context.Context, max uint, buffer int) *TailWriter {
	w := &TailWriter{
		max: max,
	}
	w.Buffer = NewBuffer(ctx, buffer, w.ondata, nil, nil)
	return w
}
