package sampler

import (
	"sync"
	"testing"
)

func TestZeroRateKeepsAll(t *testing.T) {
	s := New(0)
	for i := 0; i < 20; i++ {
		if !s.Keep() {
			t.Fatalf("rate=0: expected Keep()=true on call %d", i)
		}
	}
}

func TestNegativeRateTreatedAsZero(t *testing.T) {
	s := New(-5)
	for i := 0; i < 10; i++ {
		if !s.Keep() {
			t.Fatalf("rate=-5: expected Keep()=true on call %d", i)
		}
	}
}

func TestRateOneKeepsAll(t *testing.T) {
	s := New(1)
	for i := 0; i < 10; i++ {
		if !s.Keep() {
			t.Fatalf("rate=1: expected every line kept, failed on call %d", i)
		}
	}
}

func TestRateTwoKeepsEverySecond(t *testing.T) {
	s := New(2)
	expected := []bool{true, false, true, false, true, false}
	for i, want := range expected {
		got := s.Keep()
		if got != want {
			t.Fatalf("rate=2 call %d: got %v want %v", i, got, want)
		}
	}
}

func TestRateFiveKeepsEveryFifth(t *testing.T) {
	s := New(5)
	kept := 0
	for i := 0; i < 50; i++ {
		if s.Keep() {
			kept++
		}
	}
	if kept != 10 {
		t.Fatalf("rate=5 over 50 calls: got %d kept, want 10", kept)
	}
}

func TestResetRestartsCounting(t *testing.T) {
	s := New(3)
	s.Keep() // 1 -> kept
	s.Keep() // 2 -> skip
	s.Reset()
	if !s.Keep() {
		t.Fatal("after Reset first call should be kept")
	}
}

func TestSetRateUpdatesLive(t *testing.T) {
	s := New(0)
	s.SetRate(2)
	if !s.Keep() {
		t.Fatal("first call after SetRate(2) should be kept")
	}
	if s.Keep() {
		t.Fatal("second call after SetRate(2) should be skipped")
	}
}

func TestKeepIsConcurrencySafe(t *testing.T) {
	s := New(3)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Keep()
		}()
	}
	wg.Wait()
}
