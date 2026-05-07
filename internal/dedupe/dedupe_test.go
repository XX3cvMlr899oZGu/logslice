package dedupe

import (
	"fmt"
	"sync"
	"testing"
)

func TestZeroWindowNeverDedupes(t *testing.T) {
	f := New(0)
	for i := 0; i < 5; i++ {
		if f.IsDuplicate("same line") {
			t.Fatal("zero-window filter should never report duplicate")
		}
	}
}

func TestNegativeWindowTreatedAsZero(t *testing.T) {
	f := New(-10)
	if f.IsDuplicate("x") {
		t.Fatal("negative window should behave like zero")
	}
}

func TestFirstOccurrenceNotDuplicate(t *testing.T) {
	f := New(8)
	if f.IsDuplicate("hello") {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestSecondOccurrenceIsDuplicate(t *testing.T) {
	f := New(8)
	f.IsDuplicate("hello")
	if !f.IsDuplicate("hello") {
		t.Fatal("second occurrence should be a duplicate")
	}
	if f.Skipped != 1 {
		t.Fatalf("expected Skipped=1, got %d", f.Skipped)
	}
}

func TestWindowEviction(t *testing.T) {
	// Window of 2: after filling with a,b the entry a should be evicted.
	f := New(2)
	f.IsDuplicate("a")
	f.IsDuplicate("b")
	f.IsDuplicate("c") // evicts "a"
	if f.IsDuplicate("a") {
		t.Fatal("'a' should have been evicted from the window")
	}
}

func TestResetClearsState(t *testing.T) {
	f := New(4)
	f.IsDuplicate("line1")
	f.IsDuplicate("line1") // skipped
	f.Reset()
	if f.Skipped != 0 {
		t.Fatalf("expected Skipped=0 after Reset, got %d", f.Skipped)
	}
	if f.IsDuplicate("line1") {
		t.Fatal("after Reset, line1 should not be considered duplicate")
	}
}

func TestIsDuplicateConcurrencySafe(t *testing.T) {
	f := New(32)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			f.IsDuplicate(fmt.Sprintf("line-%d", n%10))
		}(i)
	}
	wg.Wait() // must not race
}
