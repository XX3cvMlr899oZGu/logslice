package reader

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// buildLogFile writes N log lines with timestamps spaced 1 minute apart
// starting from base, returning the file path.
func buildLogFile(t *testing.T, base time.Time, n int) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.log")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()
	for i := 0; i < n; i++ {
		ts := base.Add(time.Duration(i) * time.Minute)
		fmt.Fprintf(f, "%s INFO message number %d\n", ts.UTC().Format(time.RFC3339), i)
	}
	return path
}

func TestSeekToTimeFindsCorrectOffset(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	path := buildLogFile(t, base, 100)

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	target := base.Add(50 * time.Minute)
	offset, err := SeekToTime(f, target)
	if err != nil {
		t.Fatalf("SeekToTime error: %v", err)
	}

	if _, err := f.Seek(offset, 0); err != nil {
		t.Fatalf("seek: %v", err)
	}
	r := NewLineReaderFromFile(f)
	if !r.Next() {
		t.Fatal("expected at least one line after seek")
	}
	line := r.Line()
	if len(line) < 20 {
		t.Fatalf("line too short: %q", line)
	}
	expectedPrefix := target.UTC().Format(time.RFC3339)
	if line[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("expected line starting with %q, got %q", expectedPrefix, line[:len(expectedPrefix)])
	}
}

func TestSeekToTimeBeforeAll(t *testing.T) {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	path := buildLogFile(t, base, 20)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	offset, err := SeekToTime(f, base.Add(-1*time.Hour))
	if err != nil {
		t.Fatalf("SeekToTime error: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected offset 0 for target before all lines, got %d", offset)
	}
}

func TestSeekToTimeAfterAll(t *testing.T) {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	n := 20
	path := buildLogFile(t, base, n)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	// Target is after the last log line.
	target := base.Add(time.Duration(n+10) * time.Minute)
	offset, err := SeekToTime(f, target)
	if err != nil {
		t.Fatalf("SeekToTime error: %v", err)
	}

	// Seeking past all lines should return an offset at or beyond the last line,
	// meaning no lines are readable at that position.
	if _, err := f.Seek(offset, 0); err != nil {
		t.Fatalf("seek: %v", err)
	}
	r := NewLineReaderFromFile(f)
	if r.Next() {
		t.Errorf("expected no lines after seeking past all entries, got: %q", r.Line())
	}
}
