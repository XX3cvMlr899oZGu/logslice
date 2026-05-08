package slicer_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/slicer"
	"github.com/yourorg/logslice/internal/truncate"
)

func buildTruncatedLogFile(t *testing.T, lines int) string {
	t.Helper()
	f, err := os.CreateTemp("", "logslice-trunc-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < lines; i++ {
		ts := base.Add(time.Duration(i) * time.Second).Format(time.RFC3339)
		// pad each line to ~120 chars so truncation is visible
		padding := strings.Repeat("x", 80)
		fmt.Fprintf(f, "%s INFO %s message %d\n", ts, padding, i)
	}
	f.Close()
	return f.Name()
}

func TestSliceTruncatesLongLines(t *testing.T) {
	path := buildTruncatedLogFile(t, 10)
	defer os.Remove(path)

	tr, err := truncate.New(60)
	if err != nil {
		t.Fatalf("truncate.New: %v", err)
	}

	var out []string
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Second)

	err = slicer.Slice(slicer.Options{
		Path:      path,
		From:      from,
		To:        to,
		Truncator: tr,
		Out:       func(line string) error { out = append(out, line); return nil },
	})
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("expected output lines")
	}
	for _, l := range out {
		if len(l) > 60 {
			t.Errorf("line exceeds max 60 bytes (len=%d): %q", len(l), l)
		}
	}
}

func TestSliceNilTruncatorPassesThrough(t *testing.T) {
	path := buildTruncatedLogFile(t, 5)
	defer os.Remove(path)

	var out []string
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(5 * time.Second)

	err := slicer.Slice(slicer.Options{
		Path: path,
		From: from,
		To:   to,
		Out:  func(line string) error { out = append(out, line); return nil },
	})
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	for _, l := range out {
		if len(l) <= 60 {
			t.Errorf("expected long untruncated line, got len %d", len(l))
		}
	}
}
