package slicer_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/slicer"
)

// buildLogFile writes a temporary log file with one line per minute starting
// at base and returns its path. The caller is responsible for removal.
func buildLogFile(t *testing.T, base time.Time, count int) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	for i := 0; i < count; i++ {
		ts := base.Add(time.Duration(i) * time.Minute)
		fmt.Fprintf(f, "%s INFO log line %d\n", ts.UTC().Format(time.RFC3339), i)
	}
	return f.Name()
}

func TestSliceExtractsRange(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	path := buildLogFile(t, base, 10) // lines at 12:00 … 12:09

	opts := slicer.Options{
		From: base.Add(2 * time.Minute), // 12:02
		To:   base.Add(5 * time.Minute), // 12:05
	}

	var buf bytes.Buffer
	res, err := slicer.Slice(path, opts, &buf)
	if err != nil {
		t.Fatalf("Slice returned error: %v", err)
	}

	if res.LinesWritten != 4 {
		t.Errorf("expected 4 lines written, got %d", res.LinesWritten)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines in output, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "12:02") {
		t.Errorf("first line should be 12:02, got: %s", lines[0])
	}
	if !strings.Contains(lines[3], "12:05") {
		t.Errorf("last line should be 12:05, got: %s", lines[3])
	}
}

func TestSliceInvalidOptions(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	path := buildLogFile(t, base, 5)

	_, err := slicer.Slice(path, slicer.Options{
		From: base.Add(5 * time.Minute),
		To:   base,
	}, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error when To is before From")
	}
}

func TestSliceMissingFile(t *testing.T) {
	_, err := slicer.Slice("/nonexistent/file.log", slicer.Options{
		From: time.Now(),
		To:   time.Now().Add(time.Minute),
	}, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
