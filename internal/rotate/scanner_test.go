package rotate

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// touch creates an empty file and sets its modification time.
func touch(t *testing.T, dir, name string, modTime time.Time) string {
	t.Helper()
	p := filepath.Join(dir, name)
	f, err := os.Create(p)
	if err != nil {
		t.Fatalf("create %s: %v", p, err)
	}
	f.Close()
	if err := os.Chtimes(p, modTime, modTime); err != nil {
		t.Fatalf("chtimes %s: %v", p, err)
	}
	return p
}

func TestScanSegmentsSortedOldestFirst(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().Truncate(time.Second)

	p3 := touch(t, dir, "app.log.3", now.Add(-3*time.Hour))
	p1 := touch(t, dir, "app.log.1", now.Add(-1*time.Hour))
	p2 := touch(t, dir, "app.log.2", now.Add(-2*time.Hour))

	got, err := ScanSegments([]string{p3, p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(got))
	}
	if got[0].Path != p3 || got[1].Path != p2 || got[2].Path != p1 {
		t.Errorf("wrong order: %v", got)
	}
}

func TestScanSegmentsSkipsMissing(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().Truncate(time.Second)

	p1 := touch(t, dir, "app.log", now)
	missing := filepath.Join(dir, "ghost.log")

	got, err := ScanSegments([]string{p1, missing})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(got))
	}
	if got[0].Path != p1 {
		t.Errorf("expected %s, got %s", p1, got[0].Path)
	}
}

func TestFilterByTimeRange(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	segments := []Segment{
		{Path: "old.log", ModTime: now.Add(-5 * time.Hour)},
		{Path: "mid.log", ModTime: now.Add(-2 * time.Hour)},
		{Path: "new.log", ModTime: now},
	}

	from := now.Add(-3 * time.Hour)
	to := now.Add(-1 * time.Hour)

	got := FilterByTimeRange(segments, from, to)
	if len(got) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(got))
	}
	if got[0].Path != "mid.log" {
		t.Errorf("unexpected first segment: %s", got[0].Path)
	}
}
