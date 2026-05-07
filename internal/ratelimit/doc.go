// Package ratelimit implements a simple token-bucket rate limiter for
// controlling log line throughput during slice or tail operations.
//
// Usage:
//
//	limiter := ratelimit.New(500) // 500 lines per second
//	for _, line := range lines {
//		limiter.Wait()
//		process(line)
//	}
//
// Passing 0 as linesPerSecond disables rate limiting entirely, making
// Wait a no-op. The rate can be adjusted at runtime via SetRate.
package ratelimit
