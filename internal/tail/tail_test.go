package tail_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/tail"
)

func writeLine(t *testing.T, f *os.File, line string) {
	t.Helper()
	_, err := f.WriteString(line + "\n")
	if err != nil {
		t.Fatalf("writeLine: %v", err)
	}
}

func TestTailReceivesNewLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()

	tr, err := tail.New(path, tail.Options{PollInterval: 20 * time.Millisecond})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer tr.Stop()

	lines := tr.Lines()

	writeLine(t, f, "2024-01-01T00:00:01Z INFO hello")
	writeLine(t, f, "2024-01-01T00:00:02Z INFO world")

	received := make([]string, 0, 2)
	timeout := time.After(2 * time.Second)
	for len(received) < 2 {
		select {
		case l := <-lines:
			received = append(received, l)
		case <-timeout:
			t.Fatalf("timeout waiting for lines, got %d", len(received))
		}
	}

	if received[0] != "2024-01-01T00:00:01Z INFO hello" {
		t.Errorf("line 0 = %q", received[0])
	}
	if received[1] != "2024-01-01T00:00:02Z INFO world" {
		t.Errorf("line 1 = %q", received[1])
	}
}

func TestTailStopClosesChannel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")

	if _, err := os.Create(path); err != nil {
		t.Fatalf("create: %v", err)
	}

	tr, err := tail.New(path, tail.Options{PollInterval: 20 * time.Millisecond})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	tr.Stop()

	select {
	case _, ok := <-tr.Lines():
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for channel close")
	}
}

func TestTailMissingFile(t *testing.T) {
	_, err := tail.New("/nonexistent/path/app.log", tail.Options{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
