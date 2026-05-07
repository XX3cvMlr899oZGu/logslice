// Package ratelimit provides a token-bucket rate limiter for controlling
// the throughput of log line processing, useful when tailing or streaming
// large files without overwhelming downstream consumers.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls the rate at which lines are processed.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to linesPerSecond lines per second.
// A value of 0 disables rate limiting (unlimited throughput).
func New(linesPerSecond float64) *Limiter {
	if linesPerSecond < 0 {
		linesPerSecond = 0
	}
	return &Limiter{
		tokens:   linesPerSecond,
		max:      linesPerSecond,
		rate:     linesPerSecond,
		lastTick: time.Now(),
		clock:    time.Now,
	}
}

// Wait blocks until a token is available, then consumes one.
// If the limiter was created with linesPerSecond == 0, Wait returns immediately.
func (l *Limiter) Wait() {
	if l.rate == 0 {
		return
	}
	for {
		l.mu.Lock()
		now := l.clock()
		elapsed := now.Sub(l.lastTick).Seconds()
		l.tokens += elapsed * l.rate
		if l.tokens > l.max {
			l.tokens = l.max
		}
		l.lastTick = now
		if l.tokens >= 1 {
			l.tokens--
			l.mu.Unlock()
			return
		}
		l.mu.Unlock()
		time.Sleep(time.Duration(float64(time.Second) / l.rate))
	}
}

// SetRate updates the rate limit at runtime. A value of 0 disables limiting.
func (l *Limiter) SetRate(linesPerSecond float64) {
	if linesPerSecond < 0 {
		linesPerSecond = 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rate = linesPerSecond
	l.max = linesPerSecond
	if l.tokens > l.max {
		l.tokens = l.max
	}
}
