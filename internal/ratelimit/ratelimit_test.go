package ratelimit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewZeroRateIsUnlimited(t *testing.T) {
	l := New(0)
	// Should return immediately without blocking.
	done := make(chan struct{})
	go func() {
		l.Wait()
		close(done)
	}()
	select {
	case <-done:
		// ok
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait blocked on unlimited limiter")
	}
}

func TestNewNegativeRateTreatedAsZero(t *testing.T) {
	l := New(-5)
	if l.rate != 0 {
		t.Fatalf("expected rate 0, got %f", l.rate)
	}
}

func TestWaitConsumesToken(t *testing.T) {
	l := New(1000) // 1000 lines/sec — fast enough for test
	start := time.Now()
	for i := 0; i < 10; i++ {
		l.Wait()
	}
	// With 1000/s and 10 tokens pre-filled, should be near-instant.
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Fatalf("10 waits at 1000/s took too long: %v", elapsed)
	}
}

func TestSetRateUpdatesLimit(t *testing.T) {
	l := New(0)
	l.SetRate(500)
	if l.rate != 500 {
		t.Fatalf("expected rate 500, got %f", l.rate)
	}
}

func TestSetRateNegativeClampsToZero(t *testing.T) {
	l := New(100)
	l.SetRate(-1)
	if l.rate != 0 {
		t.Fatalf("expected rate 0 after negative SetRate, got %f", l.rate)
	}
}

func TestWaitIsConcurrencySafe(t *testing.T) {
	l := New(100000) // effectively unlimited for concurrency test
	var wg sync.WaitGroup
	var count int64
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Wait()
			atomic.AddInt64(&count, 1)
		}()
	}
	wg.Wait()
	if count != 50 {
		t.Fatalf("expected 50 completions, got %d", count)
	}
}
