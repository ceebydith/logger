package logger

import (
	"context"
	"sync"
)

// Buffer interface defines the methods for writing to the buffer and checking if it's done.
type Buffer interface {
	Write(p []byte) (int, error)
	Done() <-chan struct{}
}

// buffer struct contains the buffer's internal data and synchronization mechanisms.
type buffer struct {
	mu      sync.Mutex
	ctx     context.Context
	buffer  []byte
	ch      chan struct{}
	done    chan struct{}
	onstart func()
	onstop  func()
	ondata  func(buffer []byte)
}

// Write appends data to the buffer and signals the buffer processing routine.
func (b *buffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	select {
	case <-b.ctx.Done():
		return 0, b.ctx.Err()
	default:
		b.buffer = append(b.buffer, p...)
		go b.signal()
		return len(p), nil
	}
}

// Done returns a channel that's closed when the buffer is done processing.
func (b *buffer) Done() <-chan struct{} {
	return b.done
}

// canStop checks if the buffer can stop processing and triggers the onstop callback if necessary.
func (b *buffer) canStop() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.buffer) != 0 {
		return false
	}
	if b.onstop != nil {
		b.onstop()
	}
	close(b.ch)
	close(b.done)
	return true
}

// pop retrieves the current buffer contents and clears the buffer.
func (b *buffer) pop() []byte {
	for {
		select {
		case <-b.ch: // Clean up the go routine
			continue
		default:
			b.mu.Lock()
			defer b.mu.Unlock()
			buffer := b.buffer
			b.buffer = nil
			return buffer
		}
	}
}

// signal notifies the buffer processing routine to process the buffer.
func (b *buffer) signal() {
	defer func() {
		recover()
	}()
	b.ch <- struct{}{}
}

// run is the main processing loop that handles buffer signals and context cancellation.
func (b *buffer) run() {
	if b.onstart != nil {
		b.onstart()
	}
	for {
		select {
		case <-b.ch:
			if buffer := b.pop(); len(buffer) > 0 && b.ondata != nil {
				b.ondata(buffer)
			}
		case <-b.ctx.Done():
			if b.canStop() {
				return
			}
		}
	}
}

// NewBuffer initializes and returns a new buffer instance with the given parameters.
func NewBuffer(ctx context.Context, size int, ondata func(buffer []byte), onstart, onstop func()) *buffer {
	buf := &buffer{
		ctx:     ctx,
		ch:      make(chan struct{}, 1),
		done:    make(chan struct{}, 1),
		onstart: onstart,
		onstop:  onstop,
		ondata:  ondata,
	}
	go buf.run()
	return buf
}
