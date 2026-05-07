package slicer_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/slicer"
)

// buildDedupedLogFile writes a log file where several lines are identical.
func buildDedupedLogFile(t *testing.T) string {
	t.Helper()
	lines := []string{
		"2024-01-01T10:00:00Z INFO  starting server",
		"2024-01-01T10:00:01Z INFO  health check ok",
		"2024-01-01T10:00:02Z INFO  health check ok",
		"2024-01-01T10:00:03Z INFO  health check ok",
		"2024-01-01T10:00:04Z ERROR disk full",
		"2024-01-01T10:00:05Z INFO  health check ok",
	}
	f, err := os.CreateTemp(t.TempDir(), "dedup-*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(strings.Join(lines, "\n") + "\n")
	f.Close()
	return f.Name()
}

func TestSliceDeduplicatesLines(t *testing.T) {
	path := buildDedupedLogFile(t)
	out := strings.Builder{}

	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 0, 6, 0, time.UTC)

	err := slicer.Slice(slicer.Options{
		FilePath:     path,
		From:         from,
		To:           to,
		Writer:       &out,
		DedupeWindow: 8,
	})
	if err != nil {
		t.Fatalf("Slice returned error: %v", err)
	}

	result := out.String()
	occurrences := strings.Count(result, "health check ok")
	if occurrences != 1 {
		t.Fatalf("expected 1 unique 'health check ok' line, got %d\noutput:\n%s", occurrences, result)
	}
	if !strings.Contains(result, "disk full") {
		t.Fatal("expected 'disk full' line to be present")
	}
}

func TestSliceDedupeWindowZeroKeepsAll(t *testing.T) {
	path := buildDedupedLogFile(t)
	out := strings.Builder{}

	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 0, 6, 0, time.UTC)

	err := slicer.Slice(slicer.Options{
		FilePath:     path,
		From:         from,
		To:           to,
		Writer:       &out,
		DedupeWindow: 0,
	})
	if err != nil {
		t.Fatalf("Slice returned error: %v", err)
	}

	occurrences := strings.Count(out.String(), "health check ok")
	if occurrences != 4 {
		t.Fatalf("expected 4 'health check ok' lines with dedup disabled, got %d", occurrences)
	}
}
