// Package dedupe provides line-level deduplication for log output.
// It tracks recently seen lines using a fixed-size LRU-style ring buffer
// and suppresses consecutive or near-consecutive duplicate entries.
package dedupe

import "sync"

// Filter holds state for deduplication.
type Filter struct {
	mu      sync.Mutex
	window  []string
	size    int
	head    int
	count   int
	Skipped int64
}

// New creates a Filter that remembers the last windowSize unique lines.
// A windowSize of 0 disables deduplication (all lines pass through).
func New(windowSize int) *Filter {
	if windowSize < 0 {
		windowSize = 0
	}
	return &Filter{
		window: make([]string, windowSize),
		size:   windowSize,
	}
}

// IsDuplicate reports whether line was seen within the current window.
// If it is a duplicate, Skipped is incremented and true is returned.
// Otherwise the line is added to the window and false is returned.
func (f *Filter) IsDuplicate(line string) bool {
	if f.size == 0 {
		return false
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// Search existing entries.
	n := f.count
	if n > f.size {
		n = f.size
	}
	for i := 0; i < n; i++ {
		if f.window[i] == line {
			f.Skipped++
			return true
		}
	}

	// Evict oldest slot and store new line.
	f.window[f.head] = line
	f.head = (f.head + 1) % f.size
	f.count++
	return false
}

// Reset clears the deduplication window.
func (f *Filter) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.window = make([]string, f.size)
	f.head = 0
	f.count = 0
	f.Skipped = 0
}
