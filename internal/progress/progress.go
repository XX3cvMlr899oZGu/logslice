// Package progress provides a simple progress reporter for long-running
// log slicing operations, emitting periodic updates to an io.Writer.
package progress

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Reporter tracks bytes processed and emits progress updates.
type Reporter struct {
	mu          sync.Mutex
	total       int64
	processed   int64
	writer      io.Writer
	ticker      *time.Ticker
	done        chan struct{}
	interval    time.Duration
}

// New creates a Reporter that writes updates every interval to w.
// total is the total number of bytes expected (0 means unknown).
func New(w io.Writer, total int64, interval time.Duration) *Reporter {
	return &Reporter{
		total:    total,
		writer:   w,
		interval: interval,
		done:     make(chan struct{}),
	}
}

// Start begins emitting periodic progress lines.
func (r *Reporter) Start() {
	r.ticker = time.NewTicker(r.interval)
	go func() {
		for {
			select {
			case <-r.ticker.C:
				r.emit()
			case <-r.done:
				return
			}
		}
	}()
}

// Add records n additional bytes as processed.
func (r *Reporter) Add(n int64) {
	r.mu.Lock()
	r.processed += n
	r.mu.Unlock()
}

// Stop halts the reporter and emits a final summary line.
func (r *Reporter) Stop() {
	if r.ticker != nil {
		r.ticker.Stop()
	}
	close(r.done)
	r.emit()
}

func (r *Reporter) emit() {
	r.mu.Lock()
	processed := r.processed
	total := r.total
	r.mu.Unlock()

	if total > 0 {
		pct := float64(processed) / float64(total) * 100
		fmt.Fprintf(r.writer, "progress: %d / %d bytes (%.1f%%)\n", processed, total, pct)
	} else {
		fmt.Fprintf(r.writer, "progress: %d bytes processed\n", processed)
	}
}
