package logger

import (
	"io"
	"sync"
)

// GatedWriter is an io.Writer implementation that buffers all of
// its data into an internal buffer until it is told to let data through.
type GatedWriter struct {
	Writer io.Writer

	buf   [][]byte
	flush bool
	lock  sync.RWMutex
}

// Flush tells the GateWriter to flush any buffered data and to stop
// buffering.
func (w *GatedWriter) Flush() {
	w.lock.Lock()
	w.flush = true
	w.lock.Unlock()

	for _, p := range w.buf {
		w.Write(p)
	}
	w.buf = nil
}

// Write is
func (w *GatedWriter) Write(p []byte) (n int, err error) {
	w.lock.RLock()
	if w.flush {
		w.lock.RUnlock()
		return w.Writer.Write(p)
	}
	w.lock.RUnlock()

	w.lock.Lock()
	defer w.lock.Unlock()

	if w.flush {
		return w.Writer.Write(p)
	}

	p2 := make([]byte, len(p))
	copy(p2, p)
	w.buf = append(w.buf, p2)
	return len(p), nil
}
