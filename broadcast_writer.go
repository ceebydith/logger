package logger

import (
	"context"
	"sync"
)

// broadcastWriter struct combines the Buffer interface with broadcasting capabilities.
type broadcastWriter struct {
	Buffer
	mu        sync.RWMutex
	listeners map[chan []byte]context.Context
}

// Listen registers a new listener for the broadcastWriter.
func (w *broadcastWriter) Listen(ctx context.Context, size ...int) <-chan []byte {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.listeners == nil {
		return nil
	}

	s := 1
	if len(size) > 0 {
		s = size[0]
	}
	ch := make(chan []byte, s)
	w.listeners[ch] = ctx
	return ch
}

// onstop closes all listener channels when the buffer processing stops.
func (w *broadcastWriter) onstop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for ch := range w.listeners {
		close(ch)
	}
	w.listeners = nil
}

// ondata broadcasts data to all listeners.
func (w *broadcastWriter) ondata(buffer []byte) {
	w.delete(w.broadcast(buffer)...)
}

// broadcast sends the message to all active listeners and collects channels for deletion.
func (w *broadcastWriter) broadcast(msg []byte) []chan []byte {
	w.mu.RLock()
	defer w.mu.RUnlock()
	var dels []chan []byte
	for ch, ctx := range w.listeners {
		select {
		case ch <- msg:
		case <-ctx.Done():
			dels = append(dels, ch)
		}
	}
	return dels
}

// delete removes and closes the specified channels from the listeners map.
func (w *broadcastWriter) delete(ch ...chan []byte) {
	if len(ch) == 0 {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, c := range ch {
		if _, ok := w.listeners[c]; ok {
			close(c)
			delete(w.listeners, c)
		}
	}
}

// BroadcastWriter initializes and returns a new broadcastWriter instance with the given parameters.
func BroadcastWriter(ctx context.Context, buffer int) *broadcastWriter {
	w := &broadcastWriter{
		listeners: map[chan []byte]context.Context{},
	}
	w.Buffer = NewBuffer(ctx, buffer, w.ondata, nil, w.onstop)
	return w
}
