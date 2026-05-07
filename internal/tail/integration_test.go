package tail_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/tail"
)

// TestTailDetectsRotation verifies that the tailer resumes reading from the
// beginning of a new file after the original is replaced (log rotation).
func TestTailDetectsRotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rotate.log")

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	tr, err := tail.New(path, tail.Options{PollInterval: 20 * time.Millisecond})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer tr.Stop()

	lines := tr.Lines()

	// Write to original file.
	writeLine(t, f, "line-before-rotation")
	f.Close()

	// Simulate rotation: truncate by replacing the file.
	time.Sleep(60 * time.Millisecond)
	newF, err := os.Create(path)
	if err != nil {
		t.Fatalf("rotate create: %v", err)
	}
	defer newF.Close()
	writeLine(t, newF, "line-after-rotation")

	received := make([]string, 0, 2)
	timeout := time.After(3 * time.Second)
	for len(received) < 2 {
		select {
		case l, ok := <-lines:
			if !ok {
				t.Fatal("channel closed unexpectedly")
			}
			received = append(received, l)
		case <-timeout:
			t.Logf("received so far: %v", received)
			t.Fatal("timeout waiting for rotated lines")
		}
	}

	if received[0] != "line-before-rotation" {
		t.Errorf("expected line-before-rotation, got %q", received[0])
	}
	if received[1] != "line-after-rotation" {
		t.Errorf("expected line-after-rotation, got %q", received[1])
	}
}
