package slicer_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/slicer"
)

// TestSliceWithKeywordFilter verifies that Slice respects keyword filtering
// when a Filter is attached to the options.
func TestSliceWithKeywordFilter(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO server started",
		"2024-01-01T10:01:00Z DEBUG request received path=/health",
		"2024-01-01T10:02:00Z ERROR database timeout path=/api",
		"2024-01-01T10:03:00Z INFO request completed path=/health",
		"2024-01-01T10:04:00Z WARN disk usage high",
	}
	f := buildLogFile(t, lines)
	defer os.Remove(f)

	start, _ := time.Parse(time.RFC3339, "2024-01-01T10:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2024-01-01T10:05:00Z")

	var buf strings.Builder
	err := slicer.Slice(slicer.Options{
		FilePath:  f,
		Start:     start,
		End:       end,
		Out:       &buf,
		Keyword:   "path=",
	})
	if err != nil {
		t.Fatalf("Slice returned error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "server started") {
		t.Error("expected 'server started' to be filtered out")
	}
	if strings.Contains(output, "disk usage high") {
		t.Error("expected 'disk usage high' to be filtered out")
	}
	if !strings.Contains(output, "request received") {
		t.Error("expected 'request received' line to be present")
	}
	if !strings.Contains(output, "database timeout") {
		t.Error("expected 'database timeout' line to be present")
	}
}
