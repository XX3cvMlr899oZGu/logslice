package progress

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestEmitWithTotal(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, 1000, time.Minute)
	r.Add(500)
	r.emit()

	out := buf.String()
	if !strings.Contains(out, "500 / 1000") {
		t.Errorf("expected byte counts in output, got: %q", out)
	}
	if !strings.Contains(out, "50.0%") {
		t.Errorf("expected percentage in output, got: %q", out)
	}
}

func TestEmitWithoutTotal(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, 0, time.Minute)
	r.Add(250)
	r.emit()

	out := buf.String()
	if !strings.Contains(out, "250 bytes processed") {
		t.Errorf("expected byte count only, got: %q", out)
	}
	if strings.Contains(out, "%") {
		t.Errorf("did not expect percentage when total is 0, got: %q", out)
	}
}

func TestAddIsConcurrencySafe(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, 0, time.Minute)

	const goroutines = 50
	const addPerGoroutine = 10
	done := make(chan struct{})

	for i := 0; i < goroutines; i++ {
		go func() {
			for j := 0; j < addPerGoroutine; j++ {
				r.Add(1)
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < goroutines; i++ {
		<-done
	}

	if r.processed != goroutines*addPerGoroutine {
		t.Errorf("expected %d processed, got %d", goroutines*addPerGoroutine, r.processed)
	}
}

func TestStartAndStop(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, 100, 10*time.Millisecond)
	r.Add(40)
	r.Start()
	time.Sleep(35 * time.Millisecond)
	r.Stop()

	// At least one ticker tick plus the final Stop emit should have fired.
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) < 2 {
		t.Errorf("expected at least 2 progress lines, got %d: %q", len(lines), buf.String())
	}
}
