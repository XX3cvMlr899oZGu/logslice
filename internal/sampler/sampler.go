// Package sampler provides log line sampling by keeping every Nth line
// that passes through the pipeline. A zero or negative rate disables
// sampling so that every line is forwarded unchanged.
package sampler

import "sync/atomic"

// Sampler decides whether a given log line should be kept based on a
// 1-in-N sampling rate.
type Sampler struct {
	rate  int64
	count atomic.Int64
}

// New returns a Sampler that keeps every nth line.
// If rate <= 0 the sampler is disabled and Keep always returns true.
func New(rate int) *Sampler {
	if rate < 0 {
		rate = 0
	}
	return &Sampler{rate: int64(rate)}
}

// Keep returns true if the current line should be included in the output.
// It is safe to call from multiple goroutines.
func (s *Sampler) Keep() bool {
	if s.rate <= 0 {
		return true
	}
	n := s.count.Add(1)
	return n%s.rate == 1
}

// SetRate updates the sampling rate at runtime.
// A value <= 0 disables sampling.
func (s *Sampler) SetRate(rate int) {
	if rate < 0 {
		rate = 0
	}
	atomic.StoreInt64(&s.rate, int64(rate))
}

// Reset resets the internal line counter to zero.
func (s *Sampler) Reset() {
	s.count.Store(0)
}
