package slicer_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/aggregate"
	"github.com/yourorg/logslice/internal/slicer"
)

func buildAggregateLogFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "agg-*.log")
	if err != nil {
		t.Fatal(err)
	}
	lines := []string{
		"2024-01-01T10:00:00Z level=info  msg=started\n",
		"2024-01-01T10:01:00Z level=info  msg=running\n",
		"2024-01-01T10:02:00Z level=error msg=oops\n",
		"2024-01-01T10:03:00Z level=warn  msg=slow\n",
		"2024-01-01T10:04:00Z level=info  msg=done\n",
	}
	for _, l := range lines {
		f.WriteString(l)
	}
	f.Close()
	return f.Name()
}

func TestSliceAggregatesField(t *testing.T) {
	path := buildAggregateLogFile(t)
	defer os.Remove(path)

	counter := aggregate.New("level")

	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)

	opts := slicer.Options{
		FilePath:  path,
		From:      from,
		To:        to,
		Out:       strings.NewWriter(),
		Counter:   counter,
		FieldName: "level",
	}

	if err := slicer.Slice(opts); err != nil {
		t.Fatalf("Slice: %v", err)
	}

	if counter.Total() != 5 {
		t.Errorf("expected 5 total, got %d", counter.Total())
	}

	res := counter.Results()
	if len(res) == 0 {
		t.Fatal("expected non-empty results")
	}
	if res[0].Value != "info" || res[0].Count != 3 {
		t.Errorf("unexpected top result: %+v", res[0])
	}
}

func TestSliceNilCounterPassesThrough(t *testing.T) {
	path := buildAggregateLogFile(t)
	defer os.Remove(path)

	from := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)

	opts := slicer.Options{
		FilePath:  path,
		From:      from,
		To:        to,
		Out:       strings.NewWriter(),
		Counter:   nil,
		FieldName: "level",
	}

	if err := slicer.Slice(opts); err != nil {
		t.Fatalf("Slice with nil counter: %v", err)
	}
}
